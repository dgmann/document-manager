package opentelemetry

import (
	"go.opentelemetry.io/otel/attribute"
	"strings"
)

const (
	AppNamespace    = "documentmanager"
	RecordNamespace = "records"
)

type AttributeFn = func() attribute.KeyValue

func Attributes(applyFn ...AttributeFn) []attribute.KeyValue {
	options := make([]attribute.KeyValue, len(applyFn))
	for i, apply := range applyFn {
		options[i] = apply()
	}
	return options
}
func WithRecordId(id string) AttributeFn {
	return func() attribute.KeyValue {
		return attribute.String(createKey(AppNamespace, RecordNamespace, "id"), id)
	}
}

func WithRecordIds(ids []string) AttributeFn {
	return func() attribute.KeyValue {
		return attribute.StringSlice(createKey(AppNamespace, RecordNamespace, "id"), ids)
	}
}

func createKey(parts ...string) string {
	return strings.Join(parts, ".")
}
