package gopdf

import (
	"bufio"
	"bytes"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/jung-kurt/gofpdf"
)

type PdfCreator struct {
}

func NewPdfCreator() *PdfCreator {
	return &PdfCreator{}
}

func (creator *PdfCreator) Create(document *processor.Document) (*processor.Pdf, error) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, "mm", "A4", "")
	defer pdf.Close()
	pdf.SetTitle(document.Title, true)

	pages := getPages(document)
	for i, page := range pages {
		pageName := string(i)
		var opt gofpdf.ImageOptions
		opt.ImageType = page.Format
		opt.ReadDpi = true
		opt.AllowNegativePosition = true
		content := bytes.NewReader(page.Content)
		pdf.RegisterImageOptionsReader(pageName, opt, content)
		pdf.ImageOptions(pageName, -10, 0, -1, -1, true, opt, 0, "")
	}

	var res bytes.Buffer
	byteWriter := bufio.NewWriter(&res)
	if err := pdf.Output(byteWriter); err != nil {
		return nil, err
	}
	return &processor.Pdf{Content: res.Bytes()}, nil
}

func getPages(document *processor.Document) []*processor.Image {
	var pages []*processor.Image
	for _, category := range document.Pages {
		pages = append(pages, category)
	}
	for _, record := range document.Documents {
		subdoc := getPages(record)
		pages = append(pages, subdoc...)
	}
	return pages
}
