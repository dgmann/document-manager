package http

import (
	"github.com/dgmann/document-manager/api/internal/status"
	"net/http"
)

type HealthController struct {
	healthService *status.HealthService
}

func (h *HealthController) Status(w http.ResponseWriter, req *http.Request) {
	messages := h.healthService.Collect(req.Context())
	NewResponseWithStatus(w, messages, http.StatusOK).WriteJSON()
}
