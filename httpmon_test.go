package httpmon

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRequest_AddHeader(t *testing.T) {
	req := Request{
		Method:          GetMethod,
		HttpRequestURL:  "https://example.com",
		HttpHeader:      make(HttpHeader, 0),
		HttpRequestBody: &bytes.Buffer{},
	}
	req.AddHeader("accept", "application/json")

	assert.Equal(t, 1, len(req.HttpHeader))
}

func TestRequest_AddHeader_MultipleTimes(t *testing.T) {
	req := Request{
		Method:          GetMethod,
		HttpRequestURL:  "https://example.com",
		HttpHeader:      make(HttpHeader, 0),
		HttpRequestBody: &bytes.Buffer{},
	}
	req.AddHeader("accept", "application/json")
	req.AddHeader("accept", "application/vnd.example.com+json")

	assert.Equal(t, 1, len(req.HttpHeader))
	assert.Equal(t, HttpHeaderValues{
		"application/json",
		"application/vnd.example.com+json",
	},
		req.HttpHeader["accept"])
}

func TestRequest_Body(t *testing.T) {
	req := Request{
		Method:          GetMethod,
		HttpRequestURL:  "https://example.com",
		HttpHeader:      make(HttpHeader, 0),
		HttpRequestBody: &bytes.Buffer{},
	}

	body := []byte{
		64, 65, 66,
	}
	req.Body(bytes.NewReader(body))

	actual, _ := ioutil.ReadAll(req.HttpRequestBody)
	assert.Equal(t, body, actual)
}

func TestHttpRequestMethod_Get(t *testing.T) {
	req := GET(HttpRequestURL("https://example.com"))
	req.AddHeader("authorization", "bearer 0a1b2c3d4e5f")

	headers := make(HttpHeader, 1)
	headers["authorization"] = HttpHeaderValues{"bearer 0a1b2c3d4e5f"}

	assert.Equal(t, Method("GET"), req.requestMethod())
	assert.Equal(t, HttpRequestURL("https://example.com"), req.requestURL())
	assert.Equal(t, headers, req.requestHeader())
	assert.Nil(t, req.requestBody())
}
