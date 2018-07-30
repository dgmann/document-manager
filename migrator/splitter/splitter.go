package splitter

import (
	"io/ioutil"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"github.com/dgmann/document-manager/migrator/records/models"
	pdfModel "github.com/unidoc/unidoc/pdf/model"
	"github.com/satori/go.uuid"
	"time"
)

func Split(path string) ([]*models.SubRecord, string, error) {
	tmpDir, err := ioutil.TempDir("", "migration_")
	if err != nil {
		return nil, tmpDir, errors.Wrap(err, "error creating tmp dir")
	}
	pages, bookmarks, err := readPagesAndBookmarks(path)
	if err != nil {
		return nil, tmpDir, errors.Wrap(err, "error reading bookmarks")
	}
	splitted := splitByBookmarks(pages, bookmarks)
	subrecords, err := save(splitted, tmpDir)
	if err != nil {
		return nil, tmpDir, errors.Wrap(err, "error saving pdf")
	}
	return subrecords, tmpDir, nil
}

func readPagesAndBookmarks(input string) ([]*pdfModel.PdfPage, []*Bookmark, error) {
	f, err := os.Open(input)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	pdfReader, err := pdfModel.NewPdfReader(f)
	if err != nil {
		return nil, nil, err
	}
	bookmarks, err := getBookmarks(pdfReader)
	return pdfReader.PageList, bookmarks, err
}

type SplittedPdf struct {
	Title time.Time
	Pages []*pdfModel.PdfPage
}

func splitByBookmarks(pages []*pdfModel.PdfPage, bookmarks []*Bookmark) []*SplittedPdf {
	var pdfList = make([]*SplittedPdf, len(bookmarks))
	lastIndex := len(pages)
	if len(bookmarks) == 0 {
		return []*SplittedPdf{{Title: time.Now(), Pages: pages}}
	}
	if bookmarks[len(bookmarks)-1].PageNumber >= len(pages) {
		return []*SplittedPdf{}
	}
	for i := len(bookmarks) - 1; i >= 0; i-- {
		bookmark := bookmarks[i]
		pageRange := pages[bookmark.PageNumber:lastIndex]
		lastIndex = bookmark.PageNumber
		pdfList[i] = &SplittedPdf{Title: bookmark.Title, Pages: pageRange}
	}
	return pdfList
}

func save(pdfs []*SplittedPdf, outDir string) ([]*models.SubRecord, error) {
	var pdfList = make([]*models.SubRecord, 0, len(pdfs))
	for _, pdf := range pdfs {
		writer := pdfModel.NewPdfWriter()
		for _, p := range pdf.Pages {
			writer.AddPage(p)
		}
		p := filepath.Join(outDir, uuid.NewV4().String()+".pdf")
		f, err := os.Create(p)
		if err != nil {
			return nil, err
		}
		err = writer.Write(f)
		if err != nil {
			return nil, err
		}
		f.Close()

		pdfFile := &models.SubRecord{
			Path: p,
			Date: &pdf.Title,
		}
		pdfList = append(pdfList, pdfFile)
	}
	return pdfList, nil
}
