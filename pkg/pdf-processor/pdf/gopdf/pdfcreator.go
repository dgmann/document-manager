package gopdf

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/dgmann/document-manager/pkg/pdf-processor/processor"
	"github.com/jung-kurt/gofpdf"
)

type PdfCreator struct {
}

func NewPdfCreator() *PdfCreator {
	return &PdfCreator{}
}

type Page struct {
	Image     *processor.Image
	Bookmarks []Bookmark
}

type Bookmark struct {
	Text  string
	Level int
}

var pageSize = gofpdf.PageSizeA4

func (creator *PdfCreator) Create(document *processor.Document) (*processor.Pdf, error) {
	pdf := gofpdf.New(gofpdf.OrientationPortrait, "mm", pageSize, "")
	defer pdf.Close()
	pdf.SetTitle(document.Title, true)

	pages := getPages(document, 0)
	for i, page := range pages {
		if page.Bookmarks != nil && i%2 != 0 { // If a new document starts and we have an uneven page number add an additional empty page
			pdf.AddPage()
		}

		pageName := fmt.Sprint(i)
		var opt gofpdf.ImageOptions
		opt.ImageType = page.Image.Format
		opt.ReadDpi = true
		opt.AllowNegativePosition = true

		content := bytes.NewReader(page.Image.Content)
		info := pdf.RegisterImageOptionsReader(pageName, opt, content)
		if info.Width() > info.Height() {
			pdf.AddPageFormat(gofpdf.OrientationLandscape, pdf.GetPageSizeStr(pageSize))
		} else {
			pdf.AddPage()
		}
		if page.Bookmarks != nil {
			for _, bookmark := range page.Bookmarks {
				pdf.Bookmark(bookmark.Text, bookmark.Level, -1)
			}
		}

		width, _ := pdf.GetPageSize()
		pdf.ImageOptions(pageName, 5, 0, width-10, 0, false, opt, 0, "")
	}

	var res bytes.Buffer
	byteWriter := bufio.NewWriter(&res)
	if err := pdf.Output(byteWriter); err != nil {
		return nil, err
	}
	return &processor.Pdf{Content: res.Bytes()}, nil
}

func getPages(document *processor.Document, level int) []Page {
	var pages []Page
	for i, category := range document.Pages {
		var bookmarks []Bookmark
		if i == 0 {
			bookmarks = append(bookmarks, Bookmark{Text: document.Title, Level: level})
		}
		page := Page{Image: category, Bookmarks: bookmarks}
		pages = append(pages, page)
	}

	bmForLevelAdded := false
	for _, record := range document.Documents {
		subdoc := getPages(record, level+1)
		if len(subdoc) > 0 && !bmForLevelAdded {
			subdoc[0].Bookmarks = append([]Bookmark{{Text: document.Title, Level: level}}, subdoc[0].Bookmarks...)
			bmForLevelAdded = true
		}
		pages = append(pages, subdoc...)
	}
	return pages
}
