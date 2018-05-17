package parser

import "github.com/dgmann/document-manager/directory-watcher/models"

type Generic struct {
	Sender string
}

func (g *Generic) Parse(fileName string) *models.RecordCreate {
	return &models.RecordCreate{
		Sender: g.Sender,
	}
}
