package watcher

import (
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"github.com/dgmann/document-manager/api-client/record"
)

type Watcher interface {
	Watch(dir string, parser parser.Parser) <-chan *record.NewRecord
	Done(record *record.NewRecord)
}
