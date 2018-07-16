package validator

import (
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/pkg/errors"
	"fmt"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
)

func Validate(expected *filesystem.Index, actual *databasereader.Index, manager *shared.Manager) ([]Resolvable, *Error) {
	var err []string
	var resolvable []Resolvable

	if e := isRecordCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}
	if e := isPatientCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}
	invalidDirectories := expected.Validate()
	for _, dir := range invalidDirectories {
		err = append(err, fmt.Sprintf("Invalid directory structure: %s", dir))
	}

	resolvableInDatabase, missingInDatabase := findMissing(expected, actual.Index, filesystemErrorFactory())
	err = append(err, missingInDatabase...)
	resolvable = append(resolvable, resolvableInDatabase...)

	resolvableInFileSystem, missingInFileSystem := findMissing(actual, expected.Index, databaseErrorFactory(manager))
	err = append(err, missingInFileSystem...)
	resolvable = append(resolvable, resolvableInFileSystem...)

	return resolvable, &Error{err}
}

func databaseErrorFactory(manager *shared.Manager) func(container models.RecordContainer) Resolvable {
	return func(container models.RecordContainer) Resolvable {
		return NewDatabaseValidationError(container, manager)
	}
}

func filesystemErrorFactory() func(container models.RecordContainer) Resolvable {
	return func(container models.RecordContainer) Resolvable {
		return NewFilesystemValidationError(container)
	}
}

func findMissing(expected models.RecordIndex, actual *models.Index, errorFactory func(container models.RecordContainer) Resolvable) ([]Resolvable, []string) {
	var err []string
	var resolvable []Resolvable
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
			resolvable = append(resolvable, errorFactory(expectedRecord))
			continue
		}
		if !expectedRecord.Record().Equals(actualRecord.Record()) {
			err = append(err, fmt.Sprintf("record mismatch. Expected %s, Actual %s", expectedRecord, actualRecord))
		}
	}
	return resolvable, err
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
