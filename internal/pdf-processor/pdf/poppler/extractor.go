package poppler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/dgmann/document-manager/internal/pdf-processor/filesystem"
	"github.com/dgmann/document-manager/internal/pdf-processor/pdf"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) ToImages(ctx context.Context, data io.ReadSeeker, writer pdf.ImageSender) (int, error) {
	var errorbuf bytes.Buffer

	outdir, err := ioutil.TempDir("", "images")
	if err != nil {
		return 0, fmt.Errorf("error creating tmp dir: %w", err)
	}
	defer func() {
		if e := os.RemoveAll(outdir); e != nil && err == nil {
			err = e
		}
	}()

	cmd := exec.CommandContext(ctx, "pdfimages", "-png", "-j", "-", path.Join(outdir, "img"))
	cmd.Stdin = data
	cmd.Stderr = &errorbuf
	err = cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("error: %w. Message: %s", err, errorbuf.String())
	}

	return filesystem.ReadImagesFromDirectory(outdir, writer)
}
