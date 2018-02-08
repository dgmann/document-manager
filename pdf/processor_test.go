package pdf

import (
	"testing"
	"io"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"image/png"
	"github.com/stretchr/testify/assert"
)

type RequesterMock struct {}

func(m *RequesterMock) Do(b io.Reader) (io.ReadCloser, error) {
	result := NewResult()
	r, _ := json.Marshal(result)
	return ioutil.NopCloser(bytes.NewReader(r)), nil
}


func TestConvert(t *testing.T) {
	originalImage, err := png.Decode(bytes.NewBuffer(GetTestImage()))

	processor := PDFProcessor{requester:&RequesterMock{}}
	result, err := processor.Convert(bytes.NewBuffer(GetTestImage()))
	assert.Nil(t, err)

	assert.Equal(t, originalImage, result[0], "Images should be equal")
}
