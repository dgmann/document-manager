package splitter

import (
	"testing"
	"github.com/stretchr/testify/assert"
	pdf "github.com/unidoc/unidoc/pdf/model"
)

func TestSplitByBookmarks(t *testing.T) {
	var pages []*pdf.PdfPage
	for i := 0; i < 10; i++ {
		pages = append(pages, pdf.NewPdfPage())
	}
	bookmarks := []*Bookmark{
		{PageNumber: 0},
		{PageNumber: 1},
		{PageNumber: 5},
		{PageNumber: 9},
	}
	splitted := splitByBookmarks(pages, bookmarks)

	assert.Len(t, splitted[0].Pages, 1)
	assert.Len(t, splitted[1].Pages, 4)
	assert.Len(t, splitted[2].Pages, 4)
	assert.Len(t, splitted[3].Pages, 1)
}

func TestSplitByBookmarksNoPages(t *testing.T) {
	pages := make([]*pdf.PdfPage, 0)
	bookmarks := []*Bookmark{
		{PageNumber: 0},
		{PageNumber: 1},
		{PageNumber: 5},
		{PageNumber: 9},
	}
	splitted := splitByBookmarks(pages, bookmarks)

	assert.Len(t, splitted, 0)
}

func TestSplitByBookmarksEmptyBookmarks(t *testing.T) {
	var pages []*pdf.PdfPage
	for i := 0; i < 10; i++ {
		pages = append(pages, pdf.NewPdfPage())
	}
	var bookmarks []*Bookmark
	splitted := splitByBookmarks(pages, bookmarks)

	assert.Len(t, splitted[0].Pages, 10)
}
