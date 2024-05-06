package parser

import (
	"github.com/dgmann/document-manager/api/pkg/client"
)

type Parser interface {
	Parse(fileName string) *client.NewRecord
}
