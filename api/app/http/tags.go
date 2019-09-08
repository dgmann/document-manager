package http

import (
	"github.com/dgmann/document-manager/api/app"
	"net/http"
)

type TagController struct {
	tags app.TagService
}

func NewTagController(repository app.TagService) *TagController {
	return &TagController{tags: repository}
}

func (t *TagController) All(w http.ResponseWriter, req *http.Request) {
	tags, err := t.tags.All(req.Context())
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	NewResponseWithStatus(w, tags, http.StatusOK)
}
