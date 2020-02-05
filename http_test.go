package httpmon

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestHttpMon_Run_Error(t *testing.T) {
	httpMon := HttpMon{
		TimeOut: timeout,
		HttpClientFactory: func(o TimeOut) HttpClient {
			return &errorHttpClient{}
		},
	}
	test, err := httpMon.Run(&testHttpTestRequest{
		method: "GET",
		url:    "http://example.com",
	})

	assert.NotNil(t, err)
	assert.Nil(t, test)
}

var timeout TimeOut = TimeOut{
	Amount:   5,
	TimeUnit: Seconds,
}

type errorHttpClient struct {
}

func (ech *errorHttpClient) Run(details HttpRequestDetails) (HttpResponse, error) {
	return nil, fmt.Errorf("method: %s, url: %s", details.Method(), details.URL())
}

type testHttpTestRequest struct {
	method string
	url    string
}

func (r *testHttpTestRequest) Body() io.Reader {
	panic("implement me")
}

func (r *testHttpTestRequest) Headers() HttpHeaders {
	panic("implement me")
}

func (r *testHttpTestRequest) Method() HttpMethod {
	return HttpMethod(r.method)
}

func (r *testHttpTestRequest) URL() URL {
	return URL(r.url)
}

type successHttpClient struct {
	closeCalled bool
}

func (s *successHttpClient) Run(details HttpRequestDetails) (HttpResponse, error) {
	return &successResponse{callback: func() {
		s.closeCalled = true
	}}, nil
}

type successResponse struct {
	callback func()
}

func (s *successResponse) Close() error {
	s.callback()
	return nil
}

func (s *successResponse) StatusCode() HttpStatus {
	return 200
}

func TestHttpMon_Run(t *testing.T) {
	client := successHttpClient{closeCalled: false}
	httpMon := HttpMon{
		TimeOut: TimeOut{},
		HttpClientFactory: func(o TimeOut) HttpClient {
			return &client
		},
	}
	request := testHttpTestRequest{
		method: "GET",
		url:    "https://example.com",
	}
	test, err := httpMon.Run(&request)
	assert.Nil(t, err)
	result := test.ExpectStatus(200)
	assert.True(t, result.IsSuccess())
}
