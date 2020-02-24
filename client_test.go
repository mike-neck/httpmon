package httpmon

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpClient(t *testing.T) {
	headers := make(HttpHeader, 0)
	headers["content-type"] = HttpHeaderValues{
		"application/json",
	}
	server := httptest.NewServer(handler(200, headers, []byte(`{"name":"test"}`)))
	defer server.Close()

	var client HttpClient
	client = &DefaultHttpClient{}
	assert.NotNil(t, client)

	request := GET(HttpRequestURL(server.URL))
	request.AddHeader("accept", "application/json")

	test, err := client.Run(request)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unexpected error: %v", err))
		return
	}
	result1 := test.ExpectHeader("content-type", "application/json")
	assert.NotNil(t, result1)
	assert.True(t, result1.Success(), result1.Comparison().String())
}

func handler(status HttpResponseStatus, headers HttpHeader, body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(int(status))
		header := w.Header()
		for n, values := range headers {
			for _, v := range values {
				header.Add(string(n), string(v))
			}
		}
		_, _ = w.Write(body)
	}
}
