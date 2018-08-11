package parser

import "github.com/dgmann/document-manager/api-client/record"

type Generic struct {
	Sender string
}

func (g *Generic) Parse(fileName string) *record.NewRecord {
	return &record.NewRecord{
		Sender: g.Sender,
	}
}
