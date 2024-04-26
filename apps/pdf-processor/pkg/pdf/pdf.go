package pdf

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
)

type ImageSender interface {
	Send(*processor.Image) error
}

type ImageConverter interface {
	ToImages(ctx context.Context, data io.ReadSeeker, writer ImageSender) (int, error)
}

type Rotator interface {
	Rotate(img []byte, degrees float64) (*processor.Image, error)
}

type Creator interface {
	Create(document *processor.Document) (*processor.Pdf, error)
}

type PageCounter interface {
	Count(data io.ReadSeeker) (int, error)
}

type ConverterFactory struct {
	extractor  ImageConverter
	rasterizer ImageConverter
	counter    PageCounter
}

var ErrorExtraction = errors.New("image extraction failed")

func NewConverter(extractor ImageConverter, rasterizer ImageConverter, counter PageCounter) *ConverterFactory {
	return &ConverterFactory{extractor: extractor, rasterizer: rasterizer, counter: counter}
}

func (c *ConverterFactory) Extractor() ImageConverter {
	return &converter{
		converter: c.extractor,
		counter:   c.counter,
	}
}

func (c *ConverterFactory) Rasterizer() ImageConverter {
	return &converter{
		converter: c.rasterizer,
		counter:   c.counter,
	}
}

type converter struct {
	converter ImageConverter
	counter   PageCounter
}

func (c *converter) ToImages(ctx context.Context, data io.ReadSeeker, writer ImageSender) (int, error) {
	pageCount, err := c.counter.Count(data)
	if err != nil {
		return 0, err
	}

	_, _ = data.Seek(0, io.SeekStart)
	imagesSent, err := c.converter.ToImages(ctx, data, writer)

	if err != nil {
		return 0, fmt.Errorf("%s. %w", err, ErrorExtraction)
	}

	if imagesSent != pageCount {
		return 0, fmt.Errorf("invalid page count. Expected: %d, actual: %d. %w", pageCount, imagesSent, ErrorExtraction)
	}
	return imagesSent, nil
}
