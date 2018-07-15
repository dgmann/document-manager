package validator

import (
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/pkg/errors"
	"fmt"
	"strings"
	"github.com/dgmann/document-manager/migrator/records/models"
)

func Validate(expected *filesystem.Index, actual *databasereader.Index) *validationError {
	var err []string
	if e := isRecordCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}
	if e := isPatientCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}
	missingInDatabase := findMissing(expected, actual.Index)
	err = append(err, missingInDatabase...)

	missingInFileSystem := findMissing(actual, expected.Index)
	err = append(err, missingInFileSystem...)

	return &validationError{err}
}

func findMissing(expected models.RecordIndex, actual *models.Index) []string {
	var err []string
	for _, expectedRecord := range expected.Records() {
		patId := expectedRecord.PatientId()
		spez := expectedRecord.Spezialization()
		actualPatient, e := actual.GetPatient(patId)
		if e != nil {
			err = append(err, e.Error())
			continue
		}
		actualRecord, e := actualPatient.GetBySpezialization(spez)
		if e != nil {
			err = append(err, fmt.Sprintf("error finding matching record in %s for patient %d and spezialization %s", actual.Name, patId, spez))
			continue
		}
		if !expectedRecord.Record().Equals(actualRecord.Record()) {
			err = append(err, fmt.Sprintf("record mismatch. Expected %s, Actual %s", expectedRecord, actualRecord))
		}
	}
	return err
}

func isRecordCountEqual(expected *filesystem.Index, actual *databasereader.Index) error {
	expectedRecordCount := expected.GetTotalRecordCount()
	actualRecordCount := actual.GetTotalRecordCount()
	isRecordCountEqual := expectedRecordCount == actualRecordCount

	if !isRecordCountEqual {
		return errors.New(fmt.Sprintf("record count mismatch. Expected: %d, Actual: %d", expected.GetTotalRecordCount(), actual.GetTotalRecordCount()))
	}
	return nil
}

func isPatientCountEqual(expected *filesystem.Index, actual *databasereader.Index) error {
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
