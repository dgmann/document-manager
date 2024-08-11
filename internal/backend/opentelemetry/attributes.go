package opentelemetry

import (
	"strings"
)

const (
	AppNamespace = "documentmanager"
)

func Key(parts ...string) string {
	return strings.Join(parts, ".")
}
