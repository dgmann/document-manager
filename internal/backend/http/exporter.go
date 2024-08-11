package http

import (
	"fmt"
	"net/http"

	"github.com/dgmann/document-manager/internal/backend/datastore"

	"github.com/dgmann/document-manager/internal/backend/pdf"
)

type ExporterController struct {
	creator pdf.Creator
	records datastore.RecordService
}

func (t *ExporterController) Export(w http.ResponseWriter, req *http.Request) {
	recordsIds, ok := req.URL.Query()["id"]
	if !ok {
		NewErrorResponse(w, fmt.Errorf("please specify at least one record id"), http.StatusBadRequest).WriteJSON()
	}
	title := req.URL.Query().Get("title")
	query := datastore.NewRecordQuery(datastore.WithIds(recordsIds))
	ctx := req.Context()
	records, err := t.records.Query(ctx, query)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
	}
	res, err := t.creator.Create(req.Context(), title, records)
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	w.Header().Add("Content-Type", "application/pdf")
	NewBinaryResponseWithStatus(w, res, http.StatusOK).Write()
}
