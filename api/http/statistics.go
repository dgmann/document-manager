package http

import (
	"github.com/dgmann/document-manager/api/status"
	"net/http"
)

type StatisticsController struct {
	statisticService *status.StatisticsService
}

func (h *StatisticsController) Statistics(w http.ResponseWriter, req *http.Request) {
	messages := h.statisticService.Collect()
	NewResponseWithStatus(w, messages, http.StatusOK).WriteJSON()
}
