package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/sirupsen/logrus"
	"time"
)

func SetURL(data interface{}, baseUrl string, reader app.ModTimeReader) interface{} {
	switch data.(type) {
	case *app.Record:
		return cloneAndSetUrl(*data.(*app.Record), baseUrl, reader)
	case []app.Record:
		cloned := make([]app.Record, len(data.([]app.Record)))
		for i, m := range data.([]app.Record) {
			cloned[i] = *cloneAndSetUrl(m, baseUrl, reader)
		}
		return cloned
	}
	return data
}

func cloneAndSetUrl(record app.Record, baseUrl string, reader app.ModTimeReader) *app.Record {
	clone := record.Clone()

	setURLForRecord(&clone, baseUrl, reader)
	return &clone
}

func setURLForRecord(r *app.Record, url string, reader app.ModTimeReader) {
	if r.Tags == nil {
		r.Tags = &[]string{}
	}
	if r.Pages == nil {
		r.Pages = []app.Page{}
	}

	for i := range r.Pages {
		modified, err := reader.ModTime(app.NewKeyedGenericResource([]byte{}, r.Pages[i].Format, r.Id.Hex(), r.Pages[i].Id))
		if err != nil {
			modified = time.Now()
			logrus.Error(err)
		}

		r.Pages[i].Url = fmt.Sprintf("%s/records/%s/pages/%s?modified=%d", url, r.Id.Hex(), r.Pages[i].Id, modified.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s/archive/%s", url, r.Id.Hex())
}
