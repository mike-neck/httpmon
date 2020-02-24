package httpmon

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewHttpClient(t *testing.T) {
	var timeout = Timeout(3 * time.Second)
	var client HttpClient = NewHttpClient(timeout)
	assert.NotNil(t, client)
	assert.IsType(t, new(DefaultHttpClient), client)
}

func TestHttpClient(t *testing.T) {
	headers := make(HttpHeader, 0)
	headers["content-type"] = HttpHeaderValues{
		"application/json",
	}
	server := httptest.NewServer(handler(200, headers, []byte(`{"name":"test"}`)))
	defer server.Close()

	var client HttpClient
	client = NewHttpClient(1000000)
	assert.NotNil(t, client)

	request := GET(HttpRequestURL(server.URL))
	request.AddHeader("accept", "application/json")

	test, err := client.Run(request)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("unexpected error: %v", err))
		return
	}

	status := test.ExpectStatus(200)
	assert.NotNil(t, status)
	assert.True(t, status.Success())

	header := test.ExpectHeader("content-type", "application/json")
	assert.NotNil(t, header)
	assert.True(t, header.Success(), header.Comparison().String())
}

func TestDefaultHttpClient_Run_BuildRequestFailure(t *testing.T) {
	client := &DefaultHttpClient{
		GoRequestBuilder: func(request HttpRequest) (*http.Request, error) {
			return nil, errors.New("expected error")
		},
		GoHttpClient: nil,
	}
	request := GET("https://example.com")
	test, err := client.Run(request)
	if test != nil {
		assert.Fail(t, fmt.Sprintf("unexpected non-nil value: %v", test))
		return
	}
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "expected error")
}

func TestDefaultHttpClient_Run_GoHttpClientError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	goHttpClient := NewMockGoHttpClient(controller)

	client := &DefaultHttpClient{
		GoRequestBuilder: buildRequest,
		GoHttpClient:     goHttpClient,
	}

	goHttpClient.EXPECT().
		Run(gomock.Any()).
		Return(nil, errors.New("expected client error"))

	test, err := client.Run(GET("https://example.com/users/apps/11223344"))
	if test != nil {
		assert.Fail(t, fmt.Sprintf("unexpected non nil value: %v", test))
		return
	}
	assert.NotNil(t, err)
	assert.Equal(t, "expected client error", err.Error())
}

func handler(status HttpResponseStatus, headers HttpHeader, body []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for n, values := range headers {
			for _, v := range values {
				w.Header().Add(string(n), string(v))
			}
		}
		w.WriteHeader(int(status))
		_, _ = w.Write(body)
	}
}

func TestBuildRequest(t *testing.T) {
	request := GET("https://example.com/users/apps/200")
	request.AddHeader("accept", "application/json")
	var req *http.Request
	req, err := buildRequest(request)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("expected no error @ buildRequest: %v ", err))
		return
	}
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, "https://example.com/users/apps/200", req.URL.String())

	header := make(http.Header, 0)
	header.Add("accept", "application/json")
	assert.Equal(t, header, req.Header)
}

func TestDefaultGoHttpClient_Run(t *testing.T) {
	var d time.Duration = 1 * time.Second
	client := defaultGoHttpClient(Timeout(d))
	assert.NotNil(t, client)
	assert.IsType(t, new(DefaultGoHttpClient), client)

	headers := make(HttpHeader, 0)
	headers["content-type"] = HttpHeaderValues{
		"application/json",
	}
	server := httptest.NewServer(handler(200, headers, []byte(`{"name":"test"}`)))
	defer server.Close()

	request, _ := http.NewRequest("GET", server.URL, &bytes.Buffer{})
	response, err := client.Run(request)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("expected no error but error: %v", err))
		return
	}
	defer func() { _ = response.Body.Close() }()

	assert.Equal(t, 200, response.StatusCode)
	bs, err := ioutil.ReadAll(response.Body)
	if err != nil {
		assert.Fail(t, "unexpected error: %v", err)
		return
	}
	assert.Equal(t, []byte(`{"name":"test"}`), bs)
}

func TestDefaultGoHttpClient_Run_Failure(t *testing.T) {
	var d time.Duration = 1 * time.Millisecond
	client := defaultGoHttpClient(Timeout(d))
	assert.NotNil(t, client)
	assert.IsType(t, new(DefaultGoHttpClient), client)

	request, _ := http.NewRequest("GET", "http://localhost:8432/users/api/223344", &bytes.Buffer{})
	response, err := client.Run(request)

	if response != nil {
		assert.Fail(t, fmt.Sprintf("expected failure but succeeded: %v", *response))
		return
	}

	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "GoHttpClient#Run"))
}
