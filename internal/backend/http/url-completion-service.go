package http

import (
	"fmt"
	"net/url"

	"github.com/dgmann/document-manager/pkg/api"
)

func SetURLForRecord(record *api.Record, url url.URL) interface{} {
	return cloneAndSetUrl(*record, url)
}

func SetURLForRecordList(records []api.Record, url url.URL) interface{} {
	cloned := make([]api.Record, len(records))
	for i, m := range records {
		cloned[i] = *cloneAndSetUrl(m, url)
	}
	return cloned
}

func cloneAndSetUrl(record api.Record, url url.URL) *api.Record {
	clone := record.Clone()

	setURLForRecord(&clone, url)
	return &clone
}

func setURLForRecord(r *api.Record, url url.URL) {
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
		r.Pages = []api.Page{}
	}

	for i, page := range r.Pages {
		r.Pages[i].Url = fmt.Sprintf("%s%s/records/%s/pages/%s?modified=%d", domain, PathPrefix, r.Id, page.Id, page.UpdatedAt.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s%s/archive/%s", domain, PathPrefix, r.Id)
}
