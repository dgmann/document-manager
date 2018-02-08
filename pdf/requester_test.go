package pdf

import (
	"gopkg.in/jarcoal/httpmock.v1"
	"testing"
	"net/http"
	"bytes"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

const url  = "http://test.local"

func TestDoSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	pdfResponse := NewResult()

	responder, _ := httpmock.NewJsonResponder(200, pdfResponse)
	httpmock.RegisterResponder("POST", url, responder)

	requester := HttpRequester{
		url: url,
		client: &http.Client{},
	}

	f := bytes.NewBufferString("")
	body, err := requester.Do(f)
	assert.Nil(t, err)
	buf := new(bytes.Buffer)
	buf.ReadFrom(body)

	var resp Result
	err = json.Unmarshal(buf.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, resp, pdfResponse, "Should be equal")
}

func TestDoFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	responder, _ := httpmock.NewJsonResponder(500, `{"error"": "Test error"}`)
	httpmock.RegisterResponder("POST", url, responder)

	requester := HttpRequester{
		url: url,
		client: &http.Client{},
	}

	f := bytes.NewBufferString("")
	_, err := requester.Do(f)
	assert.NotNil(t, err)
}
