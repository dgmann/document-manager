package datastore

import (
	"testing"

	"github.com/dgmann/document-manager/pkg/api"
	"github.com/stretchr/testify/assert"
)

func TestStatus_IsNone(t *testing.T) {
	status := api.StatusNone
	assert.True(t, status.IsNone())
}

func TestStatus_IsValid_True(t *testing.T) {
	invalid := api.StatusDone
	assert.True(t, invalid.IsValid())
}

func TestStatus_IsValid_False(t *testing.T) {
	invalid := api.Status("invalid")
	assert.False(t, invalid.IsValid())
}
