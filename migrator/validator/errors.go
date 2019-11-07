package validator

import (
	"github.com/dgmann/document-manager/migrator/records/models"
	"github.com/dgmann/document-manager/migrator/shared"
	"os"
	"strings"
)

type Error struct {
	Messages []string
}

func (e *Error) Error() string {
	return strings.Join(e.Messages, ";")
}

type ResolvableValidationError interface {
	models.RecordContainer
	Resolvable
}

type Resolvable interface {
	Resolve() error
}

type DatabaseValidationError struct {
	models.RecordContainer
	manager *shared.Manager
}

func NewDatabaseValidationError(record models.RecordContainer, manager *shared.Manager) *DatabaseValidationError {
	return &DatabaseValidationError{RecordContainer: record, manager: manager}
}

func (e *DatabaseValidationError) Resolve() error {
	stmt := `delete from PdfDatabase.dbo.Befund where Id = $1`
	_, err := e.manager.Db.Exec(stmt, e.Record().Id)
	return err
}

type FilesystemValidationError struct {
	models.RecordContainer
	manager *shared.Manager
}

func NewFilesystemValidationError(record models.RecordContainer) *FilesystemValidationError {
	return &FilesystemValidationError{RecordContainer: record}
}

func (e *FilesystemValidationError) Resolve() error {
	p := e.Record().Path
	return os.Remove(p)
}
