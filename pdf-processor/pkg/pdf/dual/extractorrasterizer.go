package dual

import (
	"fmt"
	"io"

	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/sirupsen/logrus"
)

type Processor struct {
	Extractor  pdf.ImageConverter
	Rasterizer pdf.ImageConverter
	Counter    pdf.PageCounter
}

func NewProcessor(extractor, rasterizer pdf.ImageConverter, counter pdf.PageCounter) *Processor {
	return &Processor{extractor, rasterizer, counter}
}

func (processor *Processor) ToImages(data io.ReadSeeker) ([]*processor.Image, error) {
	pageCount, err := processor.Counter.Count(data)
	if err != nil {
		return nil, err
	}

	_, _ = data.Seek(0, io.SeekStart)
	images, err := processor.Extractor.ToImages(data)

	logrus.Debug(err)
	if len(images) == pageCount {
		return images, nil
	}

	logrus.Infof("extracting did not work. Fallback to rasterizing. Expected Pages: %d, actual: %d. Error: %s", pageCount, len(images), err)
	_, _ = data.Seek(0, io.SeekStart)
	images, err = processor.Rasterizer.ToImages(data)
	if len(images) == pageCount {
		return images, nil
	}
	return nil, fmt.Errorf("error rasterizing images. Expected %d, Actual: %d. Error: %s", pageCount, len(images), err)
}
