package poppler

import (
	"bytes"
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/filesystem"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) ToImages(data io.Reader) ([]*processor.Image, error) {
	var errorbuf bytes.Buffer

	outdir, err := ioutil.TempDir("", "images")
	if err != nil {
		return nil, fmt.Errorf("error creating tmp dir: %w", err)
	}
	defer func() {
		if e := os.RemoveAll(outdir); e != nil && err == nil {
			err = e
		}
	}()

	cmd := exec.Command("pdfimages", "-png", "-j", "-", path.Join(outdir, "img"))
	cmd.Stdin = data
	cmd.Stderr = &errorbuf
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error: %w. Message: %s", err, errorbuf.String())
	}
	return filesystem.ReadImagesFromDirectory(outdir)
}
