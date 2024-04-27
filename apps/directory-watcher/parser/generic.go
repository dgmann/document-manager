package parser

import (
	"github.com/dgmann/document-manager/apiclient"
)

type Generic struct {
	Sender string
}

func (g *Generic) Parse(fileName string) *apiclient.NewRecord {
	return &apiclient.NewRecord{
		Sender: g.Sender,
	}
}
