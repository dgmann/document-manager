package importer

import "fmt"

type ImportError struct {
	error
	Record *ImportableRecord
}

func NewImportError(record *ImportableRecord, err error) ImportError {
	return ImportError{err, record}
}

type ImportErrorList []ImportError

func (err ImportErrorList) Error() string {
	return fmt.Sprintf("error during import. Could not import %d records", len(err))
}
