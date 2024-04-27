package parser

import (
	"strings"
	"time"

	"github.com/dgmann/document-manager/apiclient"
)

type Fax struct {
}

func (f *Fax) Parse(fileName string) *apiclient.NewRecord {
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

	return &apiclient.NewRecord{
		Sender:     sender,
		ReceivedAt: receviedAt,
	}
}
