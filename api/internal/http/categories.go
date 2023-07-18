package http

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/dgmann/document-manager/api/internal/datastore"
	"github.com/dgmann/document-manager/api/pkg/api"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type CategoryController struct {
	categories categoryRepository
}

func (controller *CategoryController) Router() http.Handler {
	r := chi.NewRouter()
	r.Get("/", controller.All)
	r.Post("/", controller.Create)
	r.Put("/", controller.Update)
	return r
}

type categoryRepository interface {
	All(ctx context.Context) ([]api.Category, error)
	Add(ctx context.Context, category *api.Category) error
	Update(ctx context.Context, category *api.Category) error
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
	if err := controller.categories.Add(req.Context(), &body); err != nil {
		NewErrorResponse(w, err, http.StatusConflict).WriteJSON()
		return
	}
	NewResponseWithStatus(w, body, http.StatusCreated).WriteJSON()
}

func (controller *CategoryController) Update(w http.ResponseWriter, req *http.Request) {
	var body api.Category
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		NewErrorResponse(w, err, http.StatusBadRequest).WriteJSON()
		return
	}
	if err := controller.categories.Update(req.Context(), &body); err != nil {
		if errors.Is(err, &datastore.NotFoundError{}) {
			NewErrorResponse(w, err, http.StatusNotFound).WriteJSON()
			return
		}
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}
	NewResponseWithStatus(w, body, http.StatusCreated).WriteJSON()
}
