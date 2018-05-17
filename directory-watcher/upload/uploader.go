package upload

import "github.com/dgmann/document-manager/directory-watcher/models"

type Uploader interface {
	Upload(record *models.RecordCreate) error
}
