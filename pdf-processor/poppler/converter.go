package poppler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/api"
	"github.com/h2non/filetype"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

func (c *Processor) ToImages(data io.Reader) ([]*api.Image, error) {
	pdf, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	pageCount, err := getPageCount(pdf)
	if err != nil {
		return nil, errors.New("error getting page count. Message: " + err.Error())
	}

	images := make([]*api.Image, pageCount, pageCount)
	errorChan := make(chan error)
	var wg sync.WaitGroup
	for i := 0; i < pageCount; i++ {
		wg.Add(1)
		go func(pageNumber int) {
			img, err := convert(pdf, pageNumber)
			if err != nil {
				message := fmt.Sprintf("error converting page %d. Message %s", pageNumber, err.Error())
				logrus.WithError(err).Error(message)
				errorChan <- errors.New(message)
			}

			images[pageNumber] = img
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		close(errorChan)
	}()

	errString := ""
	for err := range errorChan {
		errString += "\n" + err.Error()
	}
	if errString != "" {
		e := errors.New(errString)
		logrus.WithError(e).Error("error converting pdf")
		return nil, e
	}
	return images, nil
}

func getPageCount(pdf []byte) (int, error) {
	var outbuf, grepErrorBuf, sedErrorBuf, pdfinfoErrorBuf bytes.Buffer
	input := bytes.NewBuffer(pdf)

	pdfinfo := exec.Command("pdfinfo", "-")
	pdfinfo.Stdin = input
	pdfinfo.Stderr = &pdfinfoErrorBuf

	var err error
	grep := exec.Command("grep", "Pages")
	grep.Stdin, err = pdfinfo.StdoutPipe()
	if err != nil {
		return 0, err
	}
	grep.Stderr = &grepErrorBuf

	sed := exec.Command("sed", `s/[^0-9]*//`)
	sed.Stdin, err = grep.StdoutPipe()
	if err != nil {
		return 0, err
	}
	sed.Stderr = &sedErrorBuf
	sed.Stdout = &outbuf

	if err := sed.Start(); err != nil {
		return 0, errors.New("error starting sed. Message: " + sedErrorBuf.String())
	}
	if err := grep.Start(); err != nil {
		return 0, errors.New("error starting grep. Message: " + grepErrorBuf.String())
	}
	if err := pdfinfo.Run(); err != nil {
		return 0, errors.New("error starting pdfinfo. Message: " + pdfinfoErrorBuf.String())
	}

	if err := grep.Wait(); err != nil {
		return 0, errors.New("error waiting for grep. Message: " + grepErrorBuf.String())
	}
	if err := sed.Wait(); err != nil {
		return 0, errors.New("error waiting for sed. Message: " + sedErrorBuf.String())
	}

	result := outbuf.String()
	return strconv.Atoi(strings.TrimSpace(result))
}

func convert(pdf []byte, page int) (*api.Image, error) {
	content, err := toImage(pdf, page)
	if err != nil {
		return nil, err
	}

	kind, err := filetype.Match(content)
	if err != nil {
		return nil, err
	}

	return &api.Image{Content: content, Format: strings.Trim(kind.Extension, ".")}, nil
}

func toImage(pdf []byte, page int) ([]byte, error) {
	var outbuf, errorbuf bytes.Buffer
	input := bytes.NewBuffer(pdf)

	pageNumber := strconv.Itoa(page)
	cmd := exec.Command("pdftoppm", "-f", pageNumber, "-l", pageNumber, "-png", "-jpeg", "-r", "200")
	cmd.Stdin = input
	cmd.Stdout = &outbuf
	cmd.Stderr = &errorbuf
	err := cmd.Run()
	if err != nil {
		return nil, errors.New("Error: " + err.Error() + ". Message: " + errorbuf.String())
	}
	return outbuf.Bytes(), nil
}
