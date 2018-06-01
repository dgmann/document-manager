package parser

import (
	"github.com/dgmann/document-manager/directory-watcher/models"
	"strings"
	"time"
)

type Fax struct {
}

func (f *Fax) Parse(fileName string) *models.RecordCreate {
	sender := ""
	receviedAt := time.Now()
	toParse := strings.TrimSuffix(fileName, ".pdf")

	dateSender := strings.Split(toParse, "_Telefax.")
	if len(dateSender) >= 2 {
		sender = dateSender[1]

		result, err := time.Parse("02.01.06_15.04", dateSender[0])
		if err != nil {
			println(err)
		} else {
			receviedAt = result
		}
	}

	return &models.RecordCreate{
		Sender:     sender,
		ReceivedAt: receviedAt,
	}
}