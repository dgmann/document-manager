package models

import (
	"errors"
	"fmt"
)

type Index struct {
	Data map[int]*Patient
	Name string
}

type PatientCountable interface {
	GetTotalPatientCount() int
}

type RecordCountable interface {
	GetTotalRecordCount() int
}

type Countable interface {
	PatientCountable
	RecordCountable
}

type Savable interface {
	Save(path string) error
}

func NewIndex(name string, records []RecordContainer) *Index {
	recordsByPatient := make(map[int][]RecordContainer)
	for _, record := range records {
		recordsByPatient[record.PatientId()] = append(recordsByPatient[record.PatientId()], record)
	}

	index := make(map[int]*Patient)
	for patientId, records := range recordsByPatient {
		index[patientId] = newPatientIndex(records)
	}
	return &Index{Data: index, Name: name}
}

func (i *Index) GetPatient(id int) (*Patient, error) {
	if patient, ok := i.Data[id]; ok {
		return patient, nil
	}
	return nil, errors.New(fmt.Sprintf("could not find patient with id %d in %s", id, i.Name))
}

func (i *Index) Records() []RecordContainer {
	records := make([]RecordContainer, 0)
	for _, p := range i.Data {
		patientRecords := p.Records()
		records = append(records, patientRecords...)
	}
	return records
}

func (i *Index) GetTotalPatientCount() int {
	return len(i.Data)
}

func (i *Index) GetTotalRecordCount() int {
	return len(i.Records())
}
