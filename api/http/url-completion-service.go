package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/datastore"
	"net/url"
)

func SetURLForRecord(record *datastore.Record, url url.URL) interface{} {
	return cloneAndSetUrl(*record, url)
}

func SetURLForRecordList(records []datastore.Record, url url.URL) interface{} {
	cloned := make([]datastore.Record, len(records))
	for i, m := range records {
		cloned[i] = *cloneAndSetUrl(m, url)
	}
	return cloned
}

func cloneAndSetUrl(record datastore.Record, url url.URL) *datastore.Record {
	clone := record.Clone()

	setURLForRecord(&clone, url)
	return &clone
}

func setURLForRecord(r *datastore.Record, url url.URL) {
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

	for i, page := range r.Pages {
		r.Pages[i].Url = fmt.Sprintf("%s%s/records/%s/pages/%s?modified=%d", domain, PathPrefix, r.Id.Hex(), page.Id, page.UpdatedAt.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s%s/archive/%s", domain, PathPrefix, r.Id.Hex())
}
