package pdf

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestToImages(t *testing.T) {
	result := NewResult()
	imgs, err := result.ToImages()
	assert.Nil(t, err)
	assert.Len(t, imgs, 1)
	assert.NotNil(t, imgs[0])
}
