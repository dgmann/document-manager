package http

import (
	"github.com/dgmann/document-manager/api/internal/storage"
	"github.com/dgmann/document-manager/api/internal/storage/filesystem"
	"net/http"
)

type ArchiveController struct {
	pdfs storage.ResourceLocator
}

func (a *ArchiveController) One(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	p := a.pdfs.Locate(storage.NewLocator(filesystem.PDFFileExtension, id))
	http.ServeFile(w, req, p)
}
