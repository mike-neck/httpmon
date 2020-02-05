package httpmon

import (
	"io"
	"net/http"
)

type HttpClient interface {
	Run(builder HttpRequestDetails) (HttpResponse, error)
}

type HttpClientFactory func(out TimeOut) HttpClient

type HttpResponse interface {
	io.Closer
	StatusCode() HttpStatus
}

////

func DefaultHttpClient(out TimeOut) HttpClient {
	client := http.Client{
		Timeout: out.ToDuration(),
	}
	return &defaultHttpClient{delegate: client}
}

type defaultHttpClient struct {
	delegate http.Client
}

func (dhc *defaultHttpClient) Run(details HttpRequestDetails) (HttpResponse, error) {
	request, err := BuildRequest(details)
	if err != nil {
		return nil, err
	}
	response, err := dhc.delegate.Do(request)
	if err != nil {
		return nil, err
	}
	return &defaultHttpResponse{
		body:       response.Body,
		statusCode: response.StatusCode,
	}, nil
}

type defaultHttpResponse struct {
	body       io.ReadCloser
	statusCode int
}

func (dhr *defaultHttpResponse) Close() error {
	return dhr.body.Close()
}

func (dhr *defaultHttpResponse) StatusCode() HttpStatus {
	return HttpStatus(dhr.statusCode)
}
