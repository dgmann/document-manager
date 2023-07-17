package http

import (
	"context"
	"github.com/dgmann/document-manager/api/internal/storage"
	"net/http"
)

type ArchiveController struct {
	pdfs PdfGetter
}

type PdfGetter interface {
	Get(ctx context.Context, id string) (storage.KeyedResource, error)
}

func (a *ArchiveController) One(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "recordId")
	file, err := a.pdfs.Get(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, http.StatusNotFound)
		return
	}

	w.WriteHeader(200)
	w.Header().Add("Content-Type", "application/pdf")
	if _, err := w.Write(file.Data()); err != nil {
		w.WriteHeader(500)
	}
}
