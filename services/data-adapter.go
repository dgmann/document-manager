package services

import "github.com/dgmann/document-manager-api/models"

type RecordDataAdapter interface {
	Query(query string, value interface{}) []models.Record
	Insert(sender string, pages []models.Page) *models.Record
}
