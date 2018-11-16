package importer

type ImportError struct {
	error
	Record *ImportableRecord
}

func NewImportError(record *ImportableRecord, err error) ImportError {
	return ImportError{err, record}
}
