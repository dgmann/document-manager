package shared

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

func (i *Index) GetPatient(id int) *PatientIndex {
	return i.data[id]
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
		r, ok := record.(*Record)
		if ok {
			records = append(records, r)
		}
	}
	return records
}

func (i *PatientIndex) GetBySpezialization(spez string) *Record {
	return i.data[spez].(*Record)
}
