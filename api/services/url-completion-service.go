package services

import (
	"fmt"
	"github.com/dgmann/document-manager/api/models"
	"github.com/dgmann/document-manager/api/repositories"
	"github.com/sirupsen/logrus"
	"time"
)

type ModTimeReader interface {
	ModTime(resource repositories.KeyedResource) (time.Time, error)
}

func SetURL(data interface{}, baseUrl string, reader ModTimeReader) interface{} {
	switch data.(type) {
	case *models.Record:
		return cloneAndSetUrl(data.(*models.Record), baseUrl, reader)
	case []*models.Record:
		cloned := make([]*models.Record, len(data.([]*models.Record)))
		for i, m := range data.([]*models.Record) {
			cloned[i] = cloneAndSetUrl(m, baseUrl, reader)
		}
		return cloned
	}
	return data
}

func cloneAndSetUrl(record *models.Record, baseUrl string, reader ModTimeReader) *models.Record {
	clone := record.Clone()

	setURLForRecord(clone, baseUrl, reader)
	return clone
}

func setURLForRecord(r *models.Record, url string, reader ModTimeReader) {
	if r.Tags == nil {
		r.Tags = &[]string{}
	}
	if r.Pages == nil {
		r.Pages = []*models.Page{}
	}

	for i := range r.Pages {
		modified, err := reader.ModTime(repositories.NewKeyedGenericResource([]byte{}, r.Pages[i].Format, r.Id.Hex(), r.Pages[i].Id))
		if err != nil {
			modified = time.Now()
			logrus.Error(err)
		}

		r.Pages[i].Url = fmt.Sprintf("%s/records/%s/pages/%s?modified=%d", url, r.Id.Hex(), r.Pages[i].Id, modified.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s/archive/%s", url, r.Id.Hex())
}
