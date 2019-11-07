package models

import (
	"fmt"
)

type RecordIndex interface {
	Records() []RecordContainer
}

type PageCountable interface {
	PageCount() int
}

type RecordContainer interface {
	Record() *Record
	Spezialization() string
	PatientId() int
	PageCountable
}

type Record struct {
	Id         int    `db:"Id" json:"id"`
	Name       string `db:"Name" json:"name"`
	PatId      int    `db:"Pat_Id" json:"patientId"`
	Spez       string `db:"Category" json:"category"`
	Pages      int    `db:"Pages" json:"pages"`
	Path       string `json:"path"`
	SubRecords []SubRecordContainer
}

func (r *Record) Record() *Record {
	return r
}

func (r *Record) PageCount() int {
	return r.Pages
}

func (r *Record) PatientId() int {
	return r.PatId
}

func (r *Record) Spezialization() string {
	return r.Spez
}

func (r *Record) String() string {
	return fmt.Sprintf("%d, %s: %s", r.PatId, r.Spezialization, r.Path)
}

func (r *Record) Equals(record *Record) bool {
	return r.PatId == record.PatId &&
		r.Spez == record.Spez &&
		r.Name == record.Name
}
