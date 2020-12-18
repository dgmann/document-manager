package splitter

import (
	"fmt"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
	pdfModel "github.com/unidoc/unidoc/pdf/model"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func Split(data io.ReadSeeker, outputDir string) ([]*models.SubRecord, string, error) {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return nil, "", err
	}
	tmpDir, err := ioutil.TempDir(outputDir, "migration_")
	if err != nil {
		return nil, tmpDir, errors.Wrap(err, "error creating tmp dir")
	}
	pages, bookmarks, err := readPagesAndBookmarks(data)
	if err != nil {
		return nil, tmpDir, err
	}
	splitted := splitByBookmarks(pages, bookmarks)
	subrecords, err := save(splitted, tmpDir)
	if err != nil {
		return nil, tmpDir, errors.Wrap(err, "error saving pdf")
	}
	return subrecords, tmpDir, nil
}

func readPagesAndBookmarks(data io.ReadSeeker) ([]*pdfModel.PdfPage, []*Bookmark, error) {
	pdfReader, err := pdfModel.NewPdfReader(data)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening pdf: %w", err)
	}
	bookmarks, err := getBookmarks(pdfReader)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading bookmarks: %w", err)
	}
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
		writer.SetVersion(1, 7)
		for _, p := range pdf.Pages {
			writer.AddPage(p)
		}
		pdf.Pages = nil // Set to nil in order to prevent memory problems
		p := filepath.Join(outDir, uuid.NewV4().String()+".pdf")
		f, err := os.Create(p)
		if err != nil {
			return nil, err
		}
		err = writer.Write(f)
		f.Close()
		if err != nil {
			return nil, err
		}

		pdfFile := &models.SubRecord{
			Path: p,
			Date: &pdf.Title,
		}
		pdfList = append(pdfList, pdfFile)
	}
	return pdfList, nil
}
