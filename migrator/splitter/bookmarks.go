package splitter

import (
	"github.com/unidoc/unidoc/pdf/model"
	"time"
	"github.com/unidoc/unidoc/pdf/core"
	"sort"
)

func getBookmarks(reader *model.PdfReader) ([]*Bookmark, error) {

	outlines, _, err := reader.GetOutlinesFlattened()
	if err != nil {
		return nil, err
	}
	bookmarkMap := make(map[int64]*Bookmark)
	for _, o := range outlines {
		indirectObject, ok := o.ToPdfObject().(*core.PdfIndirectObject)
		if !ok {
			continue
		}
		dict, ok := indirectObject.PdfObject.(*core.PdfObjectDictionary)
		if !ok {
			continue
		}
		destObjectArray, ok := dict.Get("Dest").(*core.PdfObjectArray)
		if !ok {
			continue
		}
		page := findPageReference(*destObjectArray)
		if page == nil {
			continue
		}
		title := dict.Get("Title")
		titleParsed, _ := time.Parse("02.01.2006 ", title.String())
		bookmarkMap[page.ObjectNumber] = &Bookmark{Title: titleParsed, ObjectNumber: page.ObjectNumber}
	}
	for i, page := range reader.PageList {
		p := page.GetPageAsIndirectObject()
		if _, ok := bookmarkMap[p.ObjectNumber]; ok {
			bookmarkMap[p.ObjectNumber].PageNumber = i
		}
	}

	var bookmarks []*Bookmark
	for _, b := range bookmarkMap {
		bookmarks = append(bookmarks, b)
	}
	sort.Slice(bookmarks, func(i, j int) bool {
		return bookmarks[i].PageNumber < bookmarks[j].PageNumber
	})
	return bookmarks, nil
}

func findPageReference(objects core.PdfObjectArray) *core.PdfIndirectObject {
	for _, object := range objects {
		if page, ok := object.(*core.PdfIndirectObject); ok {
			return page
		}
	}
	return nil
}

type Bookmark struct {
	Title        time.Time
	PageNumber   int
	ObjectNumber int64
}
