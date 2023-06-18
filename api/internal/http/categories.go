package http

import (
	"context"
	"encoding/json"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/go-chi/chi"
	"net/http"
)

type CategoryController struct {
	categories categoryRepository
}

func (controller *CategoryController) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", controller.All)
	r.Post("/", controller.Create)
	return r
}

type categoryRepository interface {
	All(ctx context.Context) ([]api.Category, error)
	Add(ctx context.Context, id, category string) error
}

func (controller *CategoryController) All(w http.ResponseWriter, req *http.Request) {
	categories, err := controller.categories.All(req.Context())
	if err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}
	NewResponse(w, categories).WriteJSON()
}

func (controller *CategoryController) Create(w http.ResponseWriter, req *http.Request) {
	var body api.Category
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}
	if err := controller.categories.Add(req.Context(), body.Id, body.Name); err != nil {
		NewErrorResponse(w, err, http.StatusConflict).WriteJSON()
		return
	}
	NewResponseWithStatus(w, body, http.StatusCreated).WriteJSON()
}
