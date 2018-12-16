package services

import (
	"fmt"
	"github.com/dgmann/document-manager/api/models"
	"github.com/sirupsen/logrus"
	"time"
)

func SetURL(data interface{}, baseUrl string, fileInfoService FileInfoService) interface{} {
	switch data.(type) {
	case *models.Record:
		return cloneAndSetUrl(data.(*models.Record), baseUrl, fileInfoService)
	case []*models.Record:
		cloned := make([]*models.Record, len(data.([]*models.Record)))
		for i, m := range data.([]*models.Record) {
			cloned[i] = cloneAndSetUrl(m, baseUrl, fileInfoService)
		}
		return cloned
	}
	return data
}

func cloneAndSetUrl(record *models.Record, baseUrl string, fileInfoService FileInfoService) *models.Record {
	clone := record.Clone()

	setURLForRecord(clone, baseUrl, fileInfoService)
	return clone
}

func setURLForRecord(r *models.Record, url string, fileInfoService FileInfoService) {
	if r.Tags == nil {
		r.Tags = &[]string{}
	}
	if r.Pages == nil {
		r.Pages = []*models.Page{}
	}

	for i := range r.Pages {
		fileInfo, err := fileInfoService.GetFileInfo(r.Id.Hex(), r.Pages[i].Id, r.Pages[i].Format)
		var modified time.Time
		if err != nil {
			modified = time.Now()
			logrus.Error(err)
		} else {
			modified = fileInfo.ModTime()
		}

		r.Pages[i].Url = fmt.Sprintf("%s/records/%s/pages/%s?modified=%d", url, r.Id.Hex(), r.Pages[i].Id, modified.Unix())
	}
	r.ArchivedPDF = fmt.Sprintf("%s/archive/%s", url, r.Id.Hex())
}
