package pdfcpu

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dgmann/document-manager/pdf-processor/filesystem"
	gim "github.com/ozankasikci/go-image-merge"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
)

type Extractor struct {
}

func NewExtractor() *Extractor {
	return &Extractor{}
}

var configuration = &pdfcpu.Configuration{ValidationMode: pdfcpu.ValidationNone, Reader15: true}

func (m *Extractor) Count(data io.Reader) (int, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return 0, nil
	}
	seeker := bytes.NewReader(b)
	ctx, err := api.ReadContext(seeker, configuration)
	if err != nil {
		return 0, err
	}

	if err := ctx.EnsurePageCount(); err != nil {
		return 0, err
	}

	return ctx.PageCount, nil
}

func (m *Extractor) ToImages(data io.Reader) ([]*processor.Image, error) {
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
	if err := api.ExtractImages(seeker, outdir, "pdf", nil, configuration); err != nil {
		return nil, fmt.Errorf("error extracting images: %w", err)
	}
	files, err := filesystem.ReadFiles(outdir)
	if err != nil {
		return nil, err
	}
	groups, err := groupImagesByPage(files)
	if err != nil {
		return nil, err
	}
	var images []*processor.Image
	for _, group := range groups {
		if len(group) == 1 {
			img, err := filesystem.ToImage(group[0])
			if err != nil {
				return nil, err
			}
			images = append(images, img)
		} else if len(group) > 1 {
			concatenated, err := concatenateImages(group)
			if err != nil {
				return nil, err
			}
			images = append(images, &processor.Image{Content: concatenated, Format: "png"})
		}
	}
	return images, nil
}

func groupImagesByPage(files []string) ([][]string, error) {
	pages := make([][]string, len(files)) // There are be at least as many files as pages. Probably more
	for _, file := range files {
		_, basename := filepath.Split(file)
		if basename == "" {
			return nil, errors.New("invalid filename " + file)
		}
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		parts := strings.Split(name, "_")
		if len(parts) < 2 {
			return nil, errors.New("invalid filename " + file)
		}
		pageNumber, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid filename %s: %w", file, err)
		}
		index := pageNumber - 1
		pages[index] = append(pages[index], file)
	}
	return pages, nil
}

func concatenateImages(files []string) ([]byte, error) {
	grids := make([]*gim.Grid, 0, len(files))
	for _, file := range files {
		grids = append(grids, &gim.Grid{ImageFilePath: file})
	}

	img, err := gim.New(grids, 1, len(files)).Merge()
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
