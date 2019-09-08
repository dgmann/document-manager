package http

import (
	"context"
	"github.com/dgmann/document-manager/api/app"
	"net/http"
	"time"
)

type checker map[string]app.HealthChecker

type HealthController struct {
	checker checker
}

func (h *HealthController) Status(w http.ResponseWriter, req *http.Request) {
	messages := make(map[string]string)
	for key, checker := range h.checker {
		ctx, cancel := context.WithTimeout(req.Context(), 2*time.Second)
		msg, err := checker.Check(ctx)
		cancel()
		if err != nil {
			messages[key] = err.Error()
		} else {
			messages[key] = msg
		}
	}
	NewResponseWithStatus(w, messages, http.StatusOK).WriteJSON()
}
