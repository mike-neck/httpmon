package httpmon

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

func NewHttpClient(timeout Timeout) HttpClient {
	return &DefaultHttpClient{
		GoRequestBuilder: buildRequest,
		GoHttpClient:     defaultGoHttpClient(timeout),
	}
}

type DefaultHttpClient struct {
	GoRequestBuilder
	GoHttpClient
}

func (client *DefaultHttpClient) Run(request HttpRequest) (HttpTest, error) {
	req, err := client.GoRequestBuilder(request)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	response, err := client.GoHttpClient.Run(req)
	if err != nil {
		return nil, err
	}
	finish := time.Now()

	duration := finish.Sub(start)

	return &DefaultHttpTest{
		Status:       HttpResponseStatus(response.StatusCode),
		Header:       response.Header,
		ResponseTime: ResponseTime(duration),
	}, nil
}

type GoRequestBuilder func(request HttpRequest) (*http.Request, error)

func buildRequest(request HttpRequest) (*http.Request, error) {
	req, err := http.NewRequest(string(request.requestMethod()), string(request.requestURL()), request.requestBody())
	if err != nil {
		return nil, &GoStandardError{
			Tag:      fmt.Sprintf("BuildRequest(%s)", request),
			Original: err,
		}
	}
	hs := req.Header
	header := request.requestHeader()
	for n, values := range header {
		for _, v := range values {
			hs.Add(string(n), string(v))
		}
	}
	return req, err
}

type GoStandardError struct {
	Tag      string
	Original error
}

func (err *GoStandardError) Error() string {
	return fmt.Sprintf("%s: (%v)", err.Tag, err.Original)
}

func (err *GoStandardError) IsTimeout() bool {
	if gerr, ok := err.Original.(net.Error); ok {
		return gerr.Timeout()
	}
	return false
}

// GoHttpClient limits go's standard http client only to its Do function.
type GoHttpClient interface {
	Run(r *http.Request) (*http.Response, error)
}

// DefaultGoHttpClient is default implementation of GoHttpClient
type DefaultGoHttpClient struct {
	Client *http.Client
}

// Run calls Do function of go's http client.
func (client *DefaultGoHttpClient) Run(r *http.Request) (*http.Response, error) {
	response, err := client.Client.Do(r)
	if err != nil {
		return nil, &GoStandardError{
			Tag:      fmt.Sprintf("GoHttpClient#Run#%s(%s)", r.URL, r.Method),
			Original: err,
		}
	}
	return response, nil
}

// defaultGoHttpClient creates GoHttpClient with given timeout
func defaultGoHttpClient(timeout Timeout) GoHttpClient {
	client := http.Client{
		Timeout: time.Duration(timeout),
	}
	return &DefaultGoHttpClient{
		Client: &client,
	}
}
