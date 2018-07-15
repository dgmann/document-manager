package validator

import (
	"strings"
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
	"os"
)

type Error struct {
	Messages []string
}

func (e *Error) Error() string {
	return strings.Join(e.Messages, ";")
}

type validationError struct {
	record models.RecordContainer
}

type Resolvable interface {
	Resolve() error
}

type DatabaseValidationError struct {
	*validationError
	manager *shared.Manager
}

func NewDatabaseValidationError(record models.RecordContainer, manager *shared.Manager) Resolvable {
	return &DatabaseValidationError{validationError: &validationError{record: record}, manager: manager}
}

func (e *DatabaseValidationError) Resolve() error {
	stmt := `delete from PdfDatabase.dbo.Befund where Id = $1`
	_, err := e.manager.Db.Exec(stmt, e.record.Record().Id)
	return err
}

type FilesystemValidationError struct {
	*validationError
	manager *shared.Manager
}

func NewFilesystemValidationError(record models.RecordContainer) Resolvable {
	return &FilesystemValidationError{validationError: &validationError{record: record}}
}

func (e *FilesystemValidationError) Resolve() error {
	p := e.record.Record().Path
	return os.Remove(p)
}
