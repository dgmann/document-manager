package pdf

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToImages(t *testing.T) {
	result := NewResult()
	imgs := result.ToImages()
	assert.Len(t, imgs, 1)
	assert.NotNil(t, imgs[0])
}
