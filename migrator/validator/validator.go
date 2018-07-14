package validator

import (
	"github.com/dgmann/document-manager/migrator/filesystem"
	"github.com/dgmann/document-manager/migrator/databasereader"
	"github.com/pkg/errors"
	"fmt"
	"strings"
)

func Validate(expected *filesystem.Index, actual *databasereader.Index) *validationError {
	var err []string
	if e := checkRecordCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}
	if e := checkPatientCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}

	for _, expectedRecord := range expected.GetRecords() {
		patId := expectedRecord.PatId
		spez := expectedRecord.Spezialization
		actualPatient, e := actual.GetPatient(patId)
		if e != nil {
			err = append(err, e.Error())
			continue
		}
		actualRecord, e := actualPatient.GetBySpezialization(spez)
		if e != nil {
			err = append(err, fmt.Sprintf("error finding matching record in database for patient %d and spezialization %s", patId, spez))
			continue
		}
		if !expectedRecord.Equals(actualRecord) {
			err = append(err, fmt.Sprintf("record mismatch. Expected %s, Actual %s", expectedRecord, actualRecord))
		}
	}
	return &validationError{err}
}

func checkRecordCountEqual(expected *filesystem.Index, actual *databasereader.Index) error {
	expectedRecordCount := expected.GetTotalCategorizableCount()
	actualRecordCount := actual.GetTotalCategorizableCount()
	isRecordCountEqual := expectedRecordCount == actualRecordCount

	if !isRecordCountEqual {
		return errors.New(fmt.Sprintf("record count mismatch. Expected: %d, Actual: %d", expected.GetTotalCategorizableCount(), actual.GetTotalCategorizableCount()))
	}
	return nil
}

func checkPatientCountEqual(expected *filesystem.Index, actual *databasereader.Index) error {
	expectedPatientCount := expected.GetTotalPatientCount()
	actualPatientCount := actual.GetTotalPatientCount()
	isPatientCountEqual := expectedPatientCount == actualPatientCount
	if !isPatientCountEqual {
		return errors.New(fmt.Sprintf("patient count mismatch. Expected: %d, Actual: %d", expected.GetTotalPatientCount(), actual.GetTotalPatientCount()))
	}
	return nil
}

type validationError struct {
	Messages []string
}

func (e *validationError) Error() string {
	return strings.Join(e.Messages, ";")
}
