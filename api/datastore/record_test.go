package datastore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatus_IsNone(t *testing.T) {
	status := StatusNone
	assert.True(t, status.IsNone())
}

func TestStatus_IsValid_True(t *testing.T) {
	invalid := StatusDone
	assert.True(t, invalid.IsValid())
}

func TestStatus_IsValid_False(t *testing.T) {
	invalid := Status("invalid")
	assert.False(t, invalid.IsValid())
}
