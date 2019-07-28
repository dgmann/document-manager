package parser

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/dgmann/document-manager/api/client"
)

type Generic struct {
	Sender string
}

func (g *Generic) Parse(fileName string) *client.NewRecord {
	return &client.NewRecord{
		CreateRecord: app.CreateRecord{
			Sender: g.Sender,
		},
	}
}
