package parser

import (
	"log/slog"
	"strings"
	"time"

	"github.com/dgmann/document-manager/api/pkg/client"
	"github.com/dgmann/document-manager/pkg/log"
)

type Fax struct {
}

func (f *Fax) Parse(fileName string) *client.NewRecord {
	sender := ""
	receviedAt := time.Now()
	toParse := strings.TrimSuffix(fileName, ".pdf")

	dateSender := strings.Split(toParse, "_Telefax.")
	if len(dateSender) >= 2 {
		sender = dateSender[1]

		result, err := time.Parse("02.01.06_15.04", dateSender[0])
		if err != nil {
			log.Logger.Info("failed to parse filename", log.ErrAttr(err), slog.String("parser", "fax"), slog.String("file", fileName))
		} else {
			receviedAt = result
		}
	}

	return &client.NewRecord{
		Sender:     sender,
		ReceivedAt: receviedAt,
	}
}
