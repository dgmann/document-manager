package parser

import "github.com/dgmann/document-manager/pkg/api"

type Generic struct {
	Sender string
}

func (g *Generic) Parse(fileName string) *api.NewRecord {
	return &api.NewRecord{
		Sender: g.Sender,
	}
}
