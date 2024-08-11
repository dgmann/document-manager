package http

import (
	"net/http"

	"github.com/dgmann/document-manager/internal/backend/datastore"
)

type TagController struct {
	tags datastore.TagService
}

func (t *TagController) All(w http.ResponseWriter, req *http.Request) {
	tags, err := t.tags.All(req.Context())
	if err != nil {
		NewErrorResponse(w, err, http.StatusInternalServerError).WriteJSON()
		return
	}

	NewResponseWithStatus(w, tags, http.StatusOK).WriteJSON()
}
