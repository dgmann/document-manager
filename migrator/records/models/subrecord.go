package models

import "time"

type SubRecordContainer interface {
	SubRecord() *SubRecord
	PageCountable
}

type SubRecord struct {
	Id             int        `db:"Id"`
	Name           string     `db:"Name"`
	Date           *time.Time `db:"Date"`
	ReceivedAt     time.Time  `db:"Timestamp"`
	State          int        `db:"State"`
	Type           *string    `db:"Type"`
	SenderNr       *string    `db:"SenderNr"`
	PathExtension  *string    `db:"PathExtension"`
	BefundId       *int       `db:"Befund_Id"`
	Pages          int        `db:"Pages"`
	PatId          *int       `db:"Pat_Id"`
	Spezialization *string
	Path           string
}

func (r *SubRecord) SubRecord() *SubRecord {
	return r
}

func (r *SubRecord) PageCount() int {
	return r.Pages
}
