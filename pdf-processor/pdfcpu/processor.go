package pdfcpu

import (
	"bytes"
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Processor struct {
}

func NewProcessor() *Processor {
	return &Processor{}
}

func (m *Processor) ToImages(data io.Reader) (images []*processor.Image, err error) {
	b, err := ioutil.ReadAll(data)
	seeker := bytes.NewReader(b)
	outdir, err := ioutil.TempDir("", "images")
	if err != nil {
		return nil, fmt.Errorf("error creating tmp dir: %w", err)
	}
	defer func() {
		if e := os.RemoveAll(outdir); e != nil && err == nil {
			err = e
		}
	}()
	log.SetDefaultLoggers()
	if err := api.ExtractImages(seeker, outdir, nil, &pdfcpu.Configuration{ValidationMode: pdfcpu.ValidationNone}); err != nil {
		return nil, fmt.Errorf("error extracting images: %w", err)
	}
	return readImagesFromDirectory(outdir)
}

func readImagesFromDirectory(dirname string) ([]*processor.Image, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	var images []*processor.Image
	for _, f := range files {
		p := path.Join(dirname, f.Name())
		extension := path.Ext(f.Name())
		content, err := ioutil.ReadFile(p)
		if err != nil {
			return nil, err
		}
		images = append(images, &processor.Image{Content: content, Format: strings.Trim(extension, ".")})
	}
	return images, nil
}
