package watcher

import (
	"github.com/dgmann/document-manager/api/pkg/client"
	"github.com/dgmann/document-manager/directory-watcher/parser"
	"github.com/dgmann/document-manager/pkg/log"
)

var logger = log.Logger

type Watcher interface {
	Watch(dir string, parser parser.Parser) <-chan *client.NewRecord
	Done(record *client.NewRecord)
}
