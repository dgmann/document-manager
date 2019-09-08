package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

func SetURLForRecord(record *app.Record, url url.URL, reader app.ModTimeReader) interface{} {
	return cloneAndSetUrl(*record, url, reader)
}

func SetURLForRecordList(records []app.Record, url url.URL, reader app.ModTimeReader) interface{} {
	cloned := make([]app.Record, len(records))
	for i, m := range records {
		cloned[i] = *cloneAndSetUrl(m, url, reader)
	}
	return cloned
}

func cloneAndSetUrl(record app.Record, url url.URL, reader app.ModTimeReader) *app.Record {
	clone := record.Clone()

	setURLForRecord(&clone, url, reader)
	return &clone
}

func setURLForRecord(r *app.Record, url url.URL, reader app.ModTimeReader) {
	domain := ""
	if url.Host != "" {
		host := url.Host
		scheme := url.Scheme
		if scheme == "" {
			scheme = "http"
		}
		domain = fmt.Sprintf("%s://%s", scheme, host)
	}

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

		r.Pages[i].Url = fmt.Sprintf("%s/records/%s/pages/%s?modified=%d", domain, r.Id.Hex(), r.Pages[i].Id, modified.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s/archive/%s", domain, r.Id.Hex())
}
