package models

import (
	"time"
	"fmt"
)

type RecordContainer interface {
	Record() *Record
	Spezialization() string
	PatientId() int
	LoadSubRecords() error
}

type Record struct {
	Id         int    `db:"Id"`
	Name       string `db:"Name"`
	PatId      int    `db:"Pat_Id"`
	Spez       string `db:"Category"`
	Path       string
	SubRecords []*SubRecord
}

func (r *Record) Record() *Record {
	return r
}

func (r *Record) PatientId() int {
	return r.PatId
}

func (r *Record) Spezialization() string {
	return r.Spez
}

func (r *Record) LoadSubRecords() error {
	return nil
}

func (r *Record) String() string {
	return fmt.Sprintf("%d, %s: %s", r.PatId, r.Spezialization, r.Path)
}

func (r *Record) Equals(record *Record) bool {
	return r.PatId == record.PatId &&
		r.Spez == record.Spez &&
		r.Name == record.Name
}

type SubRecord struct {
	Id            int        `db:"Id"`
	Name          string     `db:"Name"`
	Date          *time.Time `db:"Date"`
	ReceivedAt    time.Time  `db:"Timestamp"`
	State         int        `db:"State"`
	Type          *string    `db:"Type"`
	SenderNr      *string    `db:"SenderNr"`
	PathExtension *string    `db:"PathExtension"`
	BefundId      *int       `db:"Befund_Id"`
	PatId         *int       `db:"Pat_Id"`
	Path          string
}
