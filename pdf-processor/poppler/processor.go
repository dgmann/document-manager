package poppler

import (
	"bytes"
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/h2non/filetype"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (c *Processor) ToImages(data io.Reader) ([]*processor.Image, error) {
	pdf, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	pageCount, err := getPageCount(pdf)
	if err != nil {
		return nil, fmt.Errorf("error getting page count: %w", err)
	}

	images := make([]*processor.Image, pageCount, pageCount)
	errorChan := make(chan error)
	var wg sync.WaitGroup
	for i := 0; i < pageCount; i++ {
		wg.Add(1)
		go func(pageNumber int) {
			img, err := convert(pdf, pageNumber+1)
			if err != nil {
				message := fmt.Errorf("error converting page %d: %w", pageNumber, err)
				logrus.WithError(err).Error(message)
				errorChan <- message
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
		return nil, fmt.Errorf("error converting pdf: %s", errString)
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
		return 0, fmt.Errorf("error starting sed. Message: %s", sedErrorBuf.String())
	}
	if err := grep.Start(); err != nil {
		return 0, fmt.Errorf("error starting grep. Message: %s", grepErrorBuf.String())
	}
	if err := pdfinfo.Run(); err != nil {
		return 0, fmt.Errorf("error starting pdfinfo. Message: %s", pdfinfoErrorBuf.String())
	}

	if err := grep.Wait(); err != nil {
		return 0, fmt.Errorf("error waiting for grep. Message: %s", grepErrorBuf.String())
	}
	if err := sed.Wait(); err != nil {
		return 0, fmt.Errorf("error waiting for sed. Message: %s", sedErrorBuf.String())
	}

	result := outbuf.String()
	return strconv.Atoi(strings.TrimSpace(result))
}

func convert(pdf []byte, page int) (*processor.Image, error) {
	content, err := toImage(pdf, page)
	if err != nil {
		return nil, err
	}

	kind, err := filetype.Match(content)
	if err != nil {
		return nil, err
	}

	return &processor.Image{Content: content, Format: strings.Trim(kind.Extension, ".")}, nil
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
		return nil, fmt.Errorf("error: %w. Message: %s", err, errorbuf.String())
	}
	return outbuf.Bytes(), nil
}
