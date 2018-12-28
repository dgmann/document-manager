package http

import (
	"context"
	"github.com/dgmann/document-manager/api/app"
	"github.com/gin-gonic/gin"
	"time"
)

type checker map[string]app.HealthChecker

type HealthController struct {
	checker checker
}

func (h *HealthController) Status(c *gin.Context) {
	messages := make(map[string]string)
	for key, checker := range h.checker {
		ctx, cancel := context.WithTimeout(c, 2*time.Second)
		msg, err := checker.Check(ctx)
		cancel()
		if err != nil {
			messages[key] = err.Error()
		} else {
			messages[key] = msg
		}
	}
	c.JSON(200, messages)
}
