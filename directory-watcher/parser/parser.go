package parser

import "github.com/dgmann/document-manager/directory-watcher/models"

type Parser interface {
	Parse(fileName string) *models.RecordCreate
}
