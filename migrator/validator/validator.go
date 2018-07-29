package validator

import (
	"github.com/dgmann/document-manager/migrator/records/filesystem"
	"github.com/dgmann/document-manager/migrator/records/databasereader"
	"github.com/pkg/errors"
	"fmt"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
	"sync"
	"runtime"
	"github.com/sirupsen/logrus"
)

func Validate(actual *filesystem.Index, expected *databasereader.Index, manager *shared.Manager) ([]Resolvable, *Error) {
	var err []string
	var resolvable []Resolvable

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
	pagecountMismatch := isPageCountEqual(expected, actual)
	err = append(err, pagecountMismatch...)

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

func isPageCountEqual(expected models.RecordIndex, actual models.PatientIndex) []string {
	workerCount := runtime.NumCPU()
	runtime.GOMAXPROCS(workerCount + 1)
	errCh := make(chan error)

	records := expected.Records()
	chunk := len(records) / workerCount

	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(start int) {
			end := start + chunk

			if end > len(records) {
				end = len(records)
			}

			for j := start; j < end; j = j + 1 {
				errCh <- comparePageCount(records[j], actual)
			}
			wg.Done()
		}(i * chunk)
	}

	var err []string
	go func() {
		for e := range errCh {
			if e != nil {
				err = append(err, e.Error())
			}
		}
	}()

	wg.Wait()
	close(errCh)
	return err
}

func comparePageCount(expectedRecord models.RecordContainer, actual models.PatientIndex) error {
	patient, e := actual.GetPatient(expectedRecord.PatientId())
	if e != nil {
		return nil
	}
	actualRecord, e := patient.GetBySpezialization(expectedRecord.Spezialization())
	if e != nil {
		return nil
	}
	expectedPageCount := expectedRecord.PageCount()
	actualPageCount := actualRecord.PageCount()
	path := getPath(actualRecord, expectedRecord)
	if expectedPageCount != actualPageCount {
		return errors.New(fmt.Sprintf("page count mismatch for %s. Expected %d, Actual %d", path, expectedPageCount, actualPageCount))
	}
	if expectedPageCount == -1 || actualPageCount == -1 {
		return errors.New(fmt.Sprintf("pdf file %s is corrupted", path))
	}
	return nil
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
