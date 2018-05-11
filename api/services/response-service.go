package services

import (
	"github.com/dgmann/document-manager/api/models"
	"os"
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
)

type ResponseService struct {
	baseUrl         string
	fileInfoService FileInfoService
}

type FileInfoService interface {
	GetFileInfo(recordId, pageId string, format string) (os.FileInfo, error)
}

func NewResponseService(baseUrl string, fileInfoService FileInfoService) *ResponseService {
	return &ResponseService{baseUrl: baseUrl, fileInfoService: fileInfoService}
}

func (r *ResponseService) NewResponse(data interface{}) interface{} {
	switch data.(type) {
	case *models.Record:
		SetURL(data.(*models.Record), r.baseUrl, r.fileInfoService)
	case []*models.Record:
		for _, m := range data.([]*models.Record) {
			SetURL(m, r.baseUrl, r.fileInfoService)
		}
	}
	return data
}

func SetURL(r *models.Record, url string, fileInfoService FileInfoService) {
	if r.Tags == nil {
		r.Tags = []string{}
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
