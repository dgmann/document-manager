package http

import (
	"github.com/dgmann/document-manager/api/app"
	"github.com/go-chi/chi"
	"github.com/mongodb/mongo-go-driver/bson"
	"net/http"
)

type PatientController struct {
	records    app.RecordService
	tags       app.TagService
	categories app.CategoryService
}

func (controller *PatientController) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/:patientId/tags", controller.Tags)
	r.Get("/:patientId/categories", controller.Categories)
	r.Get("/:patientId/records", controller.Records)
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
	records, err := controller.records.Query(req.Context(), bson.M{"patientId": id})
	if err != nil {
		NewErrorResponse(w, err, 400).WriteJSON()
		return
	}
	NewResponse(w, records).WriteJSON()
}
