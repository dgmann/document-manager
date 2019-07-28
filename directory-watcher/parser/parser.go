package parser

import (
	"github.com/dgmann/document-manager/api/client"
)

type Parser interface {
	Parse(fileName string) *client.NewRecord
}
