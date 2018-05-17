package models

import "time"

type RecordCreate struct {
	Sender       string
	ReceivedAt   time.Time
	PdfPath      string
	RetryCounter int
}
