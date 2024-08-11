package pdfcpu

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dgmann/document-manager/pkg/pdf-processor/pdf"
	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"

	"github.com/dgmann/document-manager/pkg/pdf-processor/filesystem"
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

func (m *Extractor) Count(data io.ReadSeeker) (int, error) {
	ctx, err := api.ReadContext(data, configuration)
	if err != nil {
		return 0, err
	}

	if err := ctx.EnsurePageCount(); err != nil {
		return 0, err
	}

	return ctx.PageCount, nil
}

func (m *Extractor) ToImages(ctx context.Context, data io.ReadSeeker, writer pdf.ImageSender) (res int, err error) {
	outdir, err := ioutil.TempDir("", "images")
	if err != nil {
		return 0, fmt.Errorf("error creating tmp dir: %w", err)
	}
	defer func() {
		if e := os.RemoveAll(outdir); e != nil && err == nil {
			err = e
		}
	}()
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("panic occurred in pdfcpu package")
		}
	}()
	pageCount, err := m.Count(data)
	_, _ = data.Seek(0, io.SeekStart)

	if err := api.ExtractImages(data, outdir, "pdf", nil, configuration); err != nil {
		return 0, fmt.Errorf("error extracting images: %w", err)
	}
	files, err := filesystem.ReadFiles(outdir)
	if err != nil {
		return 0, err
	}
	groups, err := groupImagesByPage(files, pageCount)
	if err != nil {
		return 0, err
	}
	imagesSent := 0
	for _, group := range groups {
		if errors.Is(ctx.Err(), context.Canceled) {
			return imagesSent, ctx.Err()
		}

		if len(group) == 1 {
			img, err := filesystem.ToImage(group[0])
			if err != nil {
				return imagesSent, err
			}
			writer.Send(img)
		} else if len(group) > 1 {
			concatenated, err := concatenateImages(group)
			if err != nil {
				return imagesSent, err
			}
			if err := writer.Send(&processor.Image{Content: concatenated, Format: "png"}); err != nil {
				return imagesSent, err
			}
		}
	}
	return imagesSent, nil
}

func groupImagesByPage(files []string, pageCount int) ([][]string, error) {
	pages := make([][]string, pageCount)
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
		if index >= len(pages) {
			return nil, fmt.Errorf("there are fewer images than pages. Num Images %d, Num Pages: %d", len(pages), pageNumber)
		}
		pages[index] = append(pages[index], file)
	}
	for i, pageGroup := range pages {
		if len(pageGroup) == 0 { //There is an empty page
			return nil, fmt.Errorf("page %d does not contain any images", i+1)
		}
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
