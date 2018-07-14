package shared

import (
	"errors"
	"fmt"
)

type Index struct {
	data map[int]*PatientIndex
}

func NewIndex(records []Categorizable) *Index {
	recordsByPatient := make(map[int][]Categorizable)
	for _, record := range records {
		recordsByPatient[record.GetPatientId()] = append(recordsByPatient[record.GetPatientId()], record)
	}

	index := make(map[int]*PatientIndex)
	for patientId, records := range recordsByPatient {
		index[patientId] = newPatientIndex(records)
	}
	return &Index{data: index}
}

func (i *Index) GetPatient(id int) (*PatientIndex, error) {
	if patient, ok := i.data[id]; ok {
		return patient, nil
	}
	return nil, errors.New(fmt.Sprintf("could not retrive patient with id %d", id))
}

func (i *Index) GetAllCategorizable() []Categorizable {
	var records []Categorizable
	for _, patient := range i.data {
		related := patient.GetAllPatientRelated()
		for _, rel := range related {
			records = append(records, rel.(Categorizable))
		}
	}
	return records
}

func (i *Index) GetTotalPatientCount() int {
	return len(i.data)
}

func (i *Index) GetTotalCategorizableCount() int {
	return len(i.GetAllCategorizable())
}

type PatientIndex struct {
	data map[string]PatientRelated
}

func newPatientIndex(records []Categorizable) *PatientIndex {
	data := make(map[string]PatientRelated)
	for _, record := range records {
		data[record.GetSpezialization()] = record
	}
	return &PatientIndex{data: data}
}

func (i *PatientIndex) GetAllPatientRelated() []PatientRelated {
	var records []PatientRelated
	for _, record := range i.data {
		records = append(records, record)
	}
	return records
}

func (i *PatientIndex) GetBySpezialization(spez string) (*Record, error) {
	if record, ok := i.data[spez]; ok {
		return record.(*Record), nil
	}
	return nil, errors.New(fmt.Sprintf("no data for spezialization %s found", spez))
}
