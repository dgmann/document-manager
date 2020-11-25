package pdfcpu

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/dgmann/document-manager/pdf-processor/filesystem"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/log"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (m *Extractor) ToImages(data io.Reader) (images []*processor.Image, err error) {
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
	return filesystem.ReadImagesFromDirectory(outdir)
}

func (m *Extractor) Count(data io.Reader) (int, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return 0, nil
	}
	seeker := bytes.NewReader(b)
	return api.PageCount(seeker, &pdfcpu.Configuration{ValidationMode: pdfcpu.ValidationNone})
}
