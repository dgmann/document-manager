package services

import "github.com/dgmann/document-manager/api/models"

type ResponseService struct {
	baseUrl string
}

func NewResponseService(baseUrl string) *ResponseService {
	return &ResponseService{baseUrl: baseUrl}
}

func (r *ResponseService) NewResponse(data interface{}) interface{} {
	switch data.(type) {
	case *models.Record:
		data.(*models.Record).SetURL(r.baseUrl)
	case []*models.Record:
		for _, m := range data.([]*models.Record) {
			m.SetURL(r.baseUrl)
		}
	}
	return data
}
