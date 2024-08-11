package http

import (
	"net/http"

	"github.com/dgmann/document-manager/internal/backend/storage"
	"github.com/dgmann/document-manager/internal/backend/storage/filesystem"
)

type ArchiveController struct {
	pdfs storage.ResourceLocator
}

func (a *ArchiveController) One(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	p := a.pdfs.Locate(storage.NewLocator(filesystem.PDFFileExtension, id))
	http.ServeFile(w, req, p)
}
