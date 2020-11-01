package dual

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

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

func (processor *Processor) ToImages(data io.Reader) ([]*processor.Image, error) {
	file, err := ioutil.TempFile("", "pdf")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(file.Name())
	if _, err := io.Copy(file, data); err != nil {
		return nil, fmt.Errorf("error writing to temp file: %w", err)
	}

	pageCount, err := processor.Counter.Count(file)
	if err != nil {
		return nil, err
	}

	_, _ = file.Seek(0, io.SeekStart)
	images, err := processor.Extractor.ToImages(file)

	if len(images) == pageCount {
		return images, nil
	}

	logrus.Infof("extracting did not work. Fallback to rasterizing. Expected Pages: %d, actual: %d", pageCount, len(images))
	_, _ = file.Seek(0, io.SeekStart)
	images, err = processor.Rasterizer.ToImages(file)
	if len(images) == pageCount {
		return images, nil
	}
	return nil, fmt.Errorf("error rasterizing images. Expected %d, Actual: %d", pageCount, len(images))
}
