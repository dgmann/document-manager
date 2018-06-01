package validator

import (
	"github.com/dgmann/document-manager/migrator/shared"
	"github.com/dgmann/document-manager/migrator/filesystem"
	"github.com/dgmann/document-manager/migrator/databasereader"
	"github.com/pkg/errors"
	"fmt"
)

func Validate(expected filesystem.Index, actual databasereader.Index) error {
	var err error
	if e := checkRecordCountEqual(expected, actual); e != nil {
		err = shared.WrapError(err, e.Error())
	}
	if e := checkPatientCountEqual(expected, actual); e != nil {
		err = shared.WrapError(err, e.Error())
	}

	for _, expectedRecord := range expected.GetRecords() {
		patId := expectedRecord.PatId
		spez := expectedRecord.Spezialization
		actualRecord := actual.GetPatient(patId).GetBySpezialization(spez)
		if expectedRecord != actualRecord {
			err = shared.WrapError(err, fmt.Sprintf("record mismatch. Expected %s, Actual %s", expectedRecord, actualRecord))
		}
	}
	return err
}

func checkRecordCountEqual(expected filesystem.Index, actual databasereader.Index) error {
	expectedRecordCount := expected.GetTotalCategorizableCount()
	actualRecordCount := actual.GetTotalCategorizableCount()
	isRecordCountEqual := expectedRecordCount == actualRecordCount

	if !isRecordCountEqual {
		return errors.New(fmt.Sprintf("record count mismatch. Expected: %d, Actual: %d", expected.GetTotalCategorizableCount(), actual.GetTotalCategorizableCount()))
	}
	return nil
}

func checkPatientCountEqual(expected filesystem.Index, actual databasereader.Index) error {
	expectedPatientCount := expected.GetTotalPatientCount()
	actualPatientCount := actual.GetTotalPatientCount()
	isPatientCountEqual := expectedPatientCount == actualPatientCount
	if !isPatientCountEqual {
		return errors.New(fmt.Sprintf("patient count mismatch. Expected: %d, Actual: %d", expected.GetTotalCategorizableCount(), actual.GetTotalCategorizableCount()))
	}
	return nil
}
