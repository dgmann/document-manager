package parser

import "github.com/dgmann/document-manager/pkg/api"

type Parser interface {
	Parse(fileName string) *api.NewRecord
}
