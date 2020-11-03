package parser

import (
	"github.com/dgmann/document-manager/apiclient"
)

type Parser interface {
	Parse(fileName string) *apiclient.NewRecord
}
