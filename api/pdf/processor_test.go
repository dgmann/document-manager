package pdf

import (
	"testing"
	"io"
	"encoding/json"
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/ioutil"
	"errors"
)

type RequesterMock struct {
	mock.Mock
}

func(m *RequesterMock) Do(b io.Reader) (io.ReadCloser, error) {
	args := m.Called(b)
	return ioutil.NopCloser(args.Get(0).(io.Reader)), args.Error(1)
}


func TestConvertSuccess(t *testing.T) {
	requester := new(RequesterMock)
	input := bytes.NewBuffer(GetTestImage())
	r, _ := json.Marshal(NewResult())
	output := bytes.NewBuffer(r)

	requester.On("Do", input).Return(output, nil)

	processor := PDFProcessor{requester:requester}
	result, err := processor.Convert(bytes.NewBuffer(GetTestImage()))
	assert.Nil(t, err)

	assert.Equal(t, bytes.NewBuffer(GetTestImage()), result[0], "Images should be equal")
}

func TestConvertError(t *testing.T) {
	requester := new(RequesterMock)
	input := bytes.NewBuffer(GetTestImage())
	output := errors.New("error")
	requester.On("Do", input).Return(bytes.NewBuffer(nil), output)

	processor := PDFProcessor{requester:requester}
	_, err := processor.Convert(bytes.NewBuffer(GetTestImage()))
	assert.Equal(t, output, err)
}
