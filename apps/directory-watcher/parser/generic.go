package parser

import (
	"github.com/dgmann/document-manager/api/pkg/client"
)

type Generic struct {
	Sender string
}

func (g *Generic) Parse(fileName string) *client.NewRecord {
	return &client.NewRecord{
		Sender: g.Sender,
	}
}
