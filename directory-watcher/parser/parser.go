package parser

import "github.com/dgmann/document-manager/api-client/record"

type Parser interface {
	Parse(fileName string) *record.NewRecord
}
