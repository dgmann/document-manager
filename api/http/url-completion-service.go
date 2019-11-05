package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"github.com/dgmann/document-manager/api/storage"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

func SetURLForRecord(record *datastore.Record, url url.URL, reader storage.ModTimeReader) interface{} {
	return cloneAndSetUrl(*record, url, reader)
}

func SetURLForRecordList(records []datastore.Record, url url.URL, reader storage.ModTimeReader) interface{} {
	cloned := make([]datastore.Record, len(records))
	for i, m := range records {
		cloned[i] = *cloneAndSetUrl(m, url, reader)
	}
	return cloned
}

func cloneAndSetUrl(record datastore.Record, url url.URL, reader storage.ModTimeReader) *datastore.Record {
	clone := record.Clone()

	setURLForRecord(&clone, url, reader)
	return &clone
}

func setURLForRecord(r *datastore.Record, url url.URL, reader storage.ModTimeReader) {
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
		r.Pages = []datastore.Page{}
	}

	for i := range r.Pages {
		modified, err := reader.ModTime(storage.NewKeyedGenericResource([]byte{}, r.Pages[i].Format, r.Id.Hex(), r.Pages[i].Id))
		if err != nil {
			modified = time.Now()
			logrus.Error(err)
		}

		r.Pages[i].Url = fmt.Sprintf("%s%s/records/%s/pages/%s?modified=%d", domain, PathPrefix, r.Id.Hex(), r.Pages[i].Id, modified.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s%s/archive/%s", domain, PathPrefix, r.Id.Hex())
}
