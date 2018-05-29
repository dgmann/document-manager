package filesystem

import (
	"github.com/dgmann/document-manager/migrator/databasereader"
	"path"
)

type FileSystem struct {
	archivePath string
}

func (f *FileSystem) GetPath(pdfFile databasereader.PdfFile) string {
	return path.Join(f.archivePath, pdfFile.GetPath())
}
