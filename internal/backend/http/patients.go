package http

import (
	"net/http"
	"net/url"

	"github.com/dgmann/document-manager/internal/backend/datastore"
	"github.com/dgmann/document-manager/internal/backend/storage"
	"github.com/go-chi/chi/v5"
)

type PatientController struct {
	records    datastore.RecordService
	tags       datastore.TagService
	categories datastore.CategoryService
	images     storage.ImageService
}

func (controller *PatientController) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/{patientId}/tags", controller.Tags)
	r.Get("/{patientId}/categories", controller.Categories)
	r.Get("/{patientId}/records", controller.Records)
	return r
}

func (controller *PatientController) Tags(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "patientId")

	tags, err := controller.tags.ByPatient(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, 404).WriteJSON()
		return
	}

	NewResponse(w, tags).WriteJSON()
}

func (controller *PatientController) Categories(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "patientId")

	categories, err := controller.categories.FindByPatient(req.Context(), id)
	if err != nil {
		NewErrorResponse(w, err, 404).WriteJSON()
		return
	}

	NewResponse(w, categories).WriteJSON()
}

func (controller *PatientController) Records(w http.ResponseWriter, req *http.Request) {
	id := URLParamFromContext(req.Context(), "patientId")
	records, err := controller.records.Query(req.Context(), datastore.NewRecordQuery(datastore.WithPatientId(id)))
	if err != nil {
		NewErrorResponse(w, err, 400).WriteJSON()
		return
	}
	withUrl := SetURLForRecordList(records, url.URL{Scheme: req.URL.Scheme, Host: req.Host})
	NewResponse(w, withUrl).WriteJSON()
}
