package validator

import (
	"context"
	"fmt"
	"github.com/dgmann/document-manager/migrator/records"
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Validate(ctx context.Context, actual *filesystem.Index, expected *databasereader.Index, manager *shared.Manager) ([]ResolvableValidationError, *Error) {
	err := make([]string, 0)
	var resolvable []ResolvableValidationError

	logrus.Info("Check is the number of records is equal")
	if e := isRecordCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}

	logrus.Info("Check is the number of patients is equal")
	if e := isPatientCountEqual(expected, actual); e != nil {
		err = append(err, e.Error())
	}

	logrus.Info("Check is the directory structure is valid")
	invalidDirectories := actual.Validate()
	for _, dir := range invalidDirectories {
		err = append(err, fmt.Sprintf("Invalid directory structure: %s", dir))
	}

	logrus.Info("Find records which are stored in the filesystem but not in the database")
	resolvableInDatabase, missingInDatabase := findMissing(expected, &actual.Index, filesystemErrorFactory())
	err = append(err, missingInDatabase...)
	resolvable = append(resolvable, resolvableInDatabase...)

	logrus.Info("Find records which are stored in the database but not in the filesystem")
	resolvableInFileSystem, missingInFileSystem := findMissing(actual, &expected.Index, databaseErrorFactory(manager))
	err = append(err, missingInFileSystem...)
	resolvable = append(resolvable, resolvableInFileSystem...)

	logrus.Info("Find records where the number of pages does not equal the information stored in the database")
	pagecountMismatch := records.Parallel(ctx, records.ToRecordChannel(expected.Records()), comparePageCount(actual))
	err = append(err, pagecountMismatch...)

	invalidSubrecords := records.Parallel(ctx, records.ToRecordChannel(actual.Records()), compareSubRecordCount())
	err = append(err, invalidSubrecords...)
	return resolvable, &Error{err}
}

func databaseErrorFactory(manager *shared.Manager) func(container models.RecordContainer) ResolvableValidationError {
	return func(container models.RecordContainer) ResolvableValidationError {
		return NewDatabaseValidationError(container, manager)
	}
}

func filesystemErrorFactory() func(container models.RecordContainer) ResolvableValidationError {
	return func(container models.RecordContainer) ResolvableValidationError {
		return NewFilesystemValidationError(container)
	}
}

func findMissing(expected models.RecordIndex, actual *models.Index, errorFactory func(container models.RecordContainer) ResolvableValidationError) ([]ResolvableValidationError, []string) {
	var err []string
	var resolvable []ResolvableValidationError
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

func isRecordCountEqual(expected models.RecordCountable, actual models.RecordCountable) error {
	expectedRecordCount := expected.GetTotalRecordCount()
	actualRecordCount := actual.GetTotalRecordCount()
	isRecordCountEqual := expectedRecordCount == actualRecordCount

	if !isRecordCountEqual {
		return errors.New(fmt.Sprintf("record count mismatch. Expected: %d, Actual: %d", expected.GetTotalRecordCount(), actual.GetTotalRecordCount()))
	}
	return nil
}

func isPatientCountEqual(expected models.PatientCountable, actual models.PatientCountable) error {
	expectedPatientCount := expected.GetTotalPatientCount()
	actualPatientCount := actual.GetTotalPatientCount()
	isPatientCountEqual := expectedPatientCount == actualPatientCount
	if !isPatientCountEqual {
		return errors.New(fmt.Sprintf("patient count mismatch. Expected: %d, Actual: %d", expected.GetTotalPatientCount(), actual.GetTotalPatientCount()))
	}
	return nil
}

func comparePageCount(actual models.PatientIndex) records.ParallelExecFunc {
	return func(record models.RecordContainer) error {
		patient, e := actual.GetPatient(record.PatientId())
		if e != nil {
			return nil
		}
		actualRecord, e := patient.GetBySpezialization(record.Spezialization())
		if e != nil {
			return nil
		}
		expectedPageCount := record.PageCount()
		actualPageCount := actualRecord.PageCount()
		path := getPath(actualRecord, record)
		if expectedPageCount != actualPageCount {
			return errors.New(fmt.Sprintf("page count mismatch for %s. Expected %d, Actual %d", path, expectedPageCount, actualPageCount))
		}
		if expectedPageCount == -1 || actualPageCount == -1 {
			return errors.New(fmt.Sprintf("pdf file %s is corrupted", path))
		}
		return nil
	}
}

func getPath(a models.RecordContainer, b models.RecordContainer) string {
	if len(a.Record().Path) > 0 {
		return a.Record().Path
	}
	if len(b.Record().Path) > 0 {
		return b.Record().Path
	}
	return ""
}

func compareSubRecordCount() records.ParallelExecFunc {
	return func(record models.RecordContainer) error {
		count := 0
		for _, subrecord := range record.Record().SubRecords {
			count += subrecord.PageCount()
		}
		if count != record.PageCount() {
			return errors.New(fmt.Sprintf("record page count of %s does not match with its subrecords. Record: %d, Subrecords: %d", record.Record().Path, record.PageCount(), count))
		}
		return nil
	}
}
