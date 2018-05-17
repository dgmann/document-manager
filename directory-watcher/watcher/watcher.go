package watcher

import (
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"github.com/dgmann/document-manager/directory-watcher/models"
)

type Watcher interface {
	Watch(dir string, parser parser.Parser) <-chan *models.RecordCreate
	Done(record *models.RecordCreate)
}
