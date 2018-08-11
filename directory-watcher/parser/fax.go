package parser

import (
	"github.com/dgmann/document-manager/api-client/record"
	"strings"
	"time"
)

type Fax struct {
}

func (f *Fax) Parse(fileName string) *record.NewRecord {
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

	return &record.NewRecord{
		Sender:     sender,
		ReceivedAt: receviedAt,
	}
}
