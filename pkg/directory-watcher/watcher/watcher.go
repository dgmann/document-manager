package watcher

import (
	"github.com/dgmann/document-manager/pkg/api"
	"github.com/dgmann/document-manager/pkg/directory-watcher/parser"
	"github.com/dgmann/document-manager/pkg/log"
)

var logger = log.Logger

type Watcher interface {
	Watch(dir string, parser parser.Parser) <-chan *api.NewRecord
	Done(record *api.NewRecord)
}
