package models

import (
	"fmt"
	"errors"
)

type PatientIndex interface {
	GetPatient(id int) (*Patient, error)
}

type Patient struct {
	Data map[string]RecordContainer
}

func newPatientIndex(records []RecordContainer) *Patient {
	data := make(map[string]RecordContainer)
	for _, record := range records {
		data[record.Spezialization()] = record
	}
	return &Patient{Data: data}
}

func (i *Patient) Records() []RecordContainer {
	var records []RecordContainer
	for _, record := range i.Data {
		records = append(records, record)
	}
	return records
}

func (i *Patient) GetBySpezialization(spez string) (RecordContainer, error) {
	if record, ok := i.Data[spez]; ok {
		return record, nil
	}
	return nil, errors.New(fmt.Sprintf("no Data for spezialization %s found", spez))
}
