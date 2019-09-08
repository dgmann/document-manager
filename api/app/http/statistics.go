package http

import (
	"fmt"
	"github.com/dgmann/document-manager/api/app"
	"net/http"
)

type providers map[string]app.StatisticProvider

type StatisticsController struct {
	providers providers
}

func (h *StatisticsController) Statistics(w http.ResponseWriter, req *http.Request) {
	messages := make(map[string]string)
	for key, provider := range h.providers {
		numberOfElements, err := provider.NumberOfElements()
		if err != nil {
			messages[key] = err.Error()
		} else {
			messages[key] = fmt.Sprintf("number of elements: %d", numberOfElements)
		}
	}
	NewResponseWithStatus(w, messages, http.StatusOK)
}
