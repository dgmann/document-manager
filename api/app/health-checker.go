package app

import "context"

type HealthChecker interface {
	Check(context.Context) (string, error)
}
