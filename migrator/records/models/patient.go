package models

import (
	"fmt"
	"errors"
)

type Patient struct {
	data map[string]RecordContainer
}

func newPatientIndex(records []RecordContainer) *Patient {
	data := make(map[string]RecordContainer)
	for _, record := range records {
		data[record.Spezialization()] = record
	}
	return &Patient{data: data}
}

func (i *Patient) Records() []RecordContainer {
	var records []RecordContainer
	for _, record := range i.data {
		records = append(records, record)
	}
	return records
}

func (i *Patient) GetBySpezialization(spez string) (RecordContainer, error) {
	if record, ok := i.data[spez]; ok {
		return record, nil
	}
	return nil, errors.New(fmt.Sprintf("no data for spezialization %s found", spez))
}
