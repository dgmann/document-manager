package http

import (
	"net/http"

	"github.com/dgmann/document-manager/internal/backend/status"
)

type StatisticsController struct {
	statisticService *status.StatisticsService
}

func (h *StatisticsController) Statistics(w http.ResponseWriter, req *http.Request) {
	messages := h.statisticService.Collect(req.Context())
	NewResponseWithStatus(w, messages, http.StatusOK).WriteJSON()
}
