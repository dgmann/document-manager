package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"net/http"
)

type ExporterController struct {
	creator app.PdfCreator
	records app.RecordService
}

func NewExporterController(creator app.PdfCreator, records app.RecordService) *ExporterController {
	return &ExporterController{creator: creator, records: records}
}

func (t *ExporterController) Export(w http.ResponseWriter, req *http.Request) {
	recordsIds, ok := req.URL.Query()["id"]
	if !ok {
		NewErrorResponse(w, fmt.Errorf("please specify at least one record id"), http.StatusBadRequest).WriteJSON()
	}
	query := app.NewRecordQuery().SetIds(recordsIds)
	ctx := req.Context()
	records, err := t.records.Query(ctx, query)
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
	}
	pdf, err := t.creator.CreatePdf(req.Context(), "Test", records)

	w.Header().Add("Content-Type", "application/pdf")
	NewResponseWithStatus(w, pdf, http.StatusOK).Write()
}
