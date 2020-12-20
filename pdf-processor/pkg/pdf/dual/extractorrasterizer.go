package dual

import (
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
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
	file, err := ioutil.TempFile("", "pdf-*.pdf")
	if err != nil {
		return nil, fmt.Errorf("error creating temp file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			logrus.Warn(err)
		}
		if err := os.Remove(file.Name()); err != nil {
			logrus.Warn(err)
		}
	}()
	if _, err := io.Copy(file, data); err != nil {
		return nil, fmt.Errorf("error writing to temp file: %w", err)
	}

	_, _ = file.Seek(0, io.SeekStart)
	pageCount, err := processor.Counter.Count(file)
	if err != nil {
		return nil, err
	}

	_, _ = file.Seek(0, io.SeekStart)
	images, err := processor.Extractor.ToImages(file)

	logrus.Debug(err)
	if len(images) == pageCount {
		return images, nil
	}

	logrus.Infof("extracting did not work. Fallback to rasterizing. Expected Pages: %d, actual: %d. Error: %s", pageCount, len(images), err)
	_, _ = file.Seek(0, io.SeekStart)
	images, err = processor.Rasterizer.ToImages(file)
	if len(images) == pageCount {
		return images, nil
	}
	return nil, fmt.Errorf("error rasterizing images. Expected %d, Actual: %d. Error: %s", pageCount, len(images), err)
}
