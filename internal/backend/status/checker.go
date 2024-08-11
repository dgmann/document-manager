package status

import (
	"context"
)

type HealthCheckable interface {
	Check(context.Context) (string, error)
}

type HealthCheckables map[string]HealthCheckable

type HealthService struct {
	healthCheckables HealthCheckables
}

func NewHealthService(checkables HealthCheckables) *HealthService {
	return &HealthService{checkables}
}

func (c *HealthService) Collect(ctx context.Context) map[string]string {
	messages := make(map[string]string)
	for key, checker := range c.healthCheckables {
		msg, err := checker.Check(ctx)
		if err != nil {
			messages[key] = err.Error()
		} else {
			messages[key] = msg
		}
	}
	return messages
}
