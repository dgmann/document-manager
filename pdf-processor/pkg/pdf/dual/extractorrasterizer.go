package dual

import (
	"bytes"
	"fmt"
	"github.com/dgmann/document-manager/pdf-processor/pkg/pdf"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
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
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	seeker := bytes.NewReader(b)

	pageCount, err := processor.Counter.Count(seeker)
	if err != nil {
		return nil, err
	}

	_, _ = seeker.Seek(0, io.SeekStart)
	images, err := processor.Extractor.ToImages(seeker)

	if len(images) == pageCount {
		return images, nil
	}

	logrus.Infof("extracting did not work. Fallback to rasterizing. Expected Pages: %d, actual: %d", pageCount, len(images))
	_, _ = seeker.Seek(0, io.SeekStart)
	images, err = processor.Rasterizer.ToImages(seeker)
	if len(images) == pageCount {
		return images, nil
	}
	return nil, fmt.Errorf("error rasterizing images. Expected %d, Actual: %d", pageCount, len(images))
}
