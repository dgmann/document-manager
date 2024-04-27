package watcher

import (
	"github.com/dgmann/document-manager/apiclient"
	"github.com/dgmann/document-manager/directory-watcher/parser"
)

type Watcher interface {
	Watch(dir string, parser parser.Parser) <-chan *apiclient.NewRecord
	Done(record *apiclient.NewRecord)
}
