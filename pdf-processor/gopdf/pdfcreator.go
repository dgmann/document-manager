package gopdf

import (
	"bufio"
	"bytes"
	"github.com/dgmann/document-manager/pdf-processor/pkg/processor"
	"github.com/jung-kurt/gofpdf"
)

type PdfCreator struct {
}

func (creator *PdfCreator) Create(document processor.Document) error {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, "mm", "A4", "")
	pdf.AddPage()
	pdf.SetTitle(document.Title, true)

	var res bytes.Buffer
	byteWriter := bufio.NewWriter(&res)
	if err := pdf.Output(byteWriter); err != nil {
		return err
	}
	return nil
}
