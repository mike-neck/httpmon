package httpmon

import (
	"fmt"
	"io"
	"time"
)

type HttpRequestMethod func(url HttpRequestURL) HttpRequest

////
var GET HttpRequestMethod = func(url HttpRequestURL) HttpRequest {
	header := make(HttpHeader, 0)
	return &Request{
		Method:          GetMethod,
		HttpRequestURL:  url,
		HttpHeader:      header,
		HttpRequestBody: nil,
	}
}

var POST HttpRequestMethod = func(url HttpRequestURL) HttpRequest {
	header := make(HttpHeader, 0)
	return &Request{
		Method:          PostMethod,
		HttpRequestURL:  url,
		HttpHeader:      header,
		HttpRequestBody: nil,
	}
}

////////

type Method string

const GetMethod Method = "GET"
const PostMethod Method = "POST"
const PutMethod Method = "PUT"
const DeleteMethod Method = "DELETE"

////////

type Request struct {
	Method
	HttpRequestURL
	HttpHeader
	HttpRequestBody
}

func (req *Request) AddHeader(name HttpHeaderName, value HttpHeaderValue) {
	if values, ok := req.HttpHeader[name]; ok {
		req.HttpHeader[name] = append(values, value)
	} else {
		headerValues := make(HttpHeaderValues, 1)
		headerValues[0] = value
		req.HttpHeader[name] = headerValues
	}
}

func (req *Request) Body(body HttpRequestBody) {
	req.HttpRequestBody = body
}

func (req *Request) requestURL() HttpRequestURL {
	return req.HttpRequestURL
}

func (req *Request) requestHeader() HttpHeader {
	return req.HttpHeader
}

func (req *Request) requestBody() HttpRequestBody {
	return req.HttpRequestBody
}

func (req *Request) requestMethod() Method {
	return req.Method
}

type HttpHeaderName string

type HttpHeaderValue string

type HttpHeaderValues []HttpHeaderValue

type HttpHeader map[HttpHeaderName]HttpHeaderValues

type HttpRequestURL string

type HttpRequestBody io.Reader

type HttpRequest interface {
	testRequest
	AddHeader(name HttpHeaderName, value HttpHeaderValue)
	Body(body HttpRequestBody)
}

type testRequest interface {
	requestMethod() Method
	requestURL() HttpRequestURL
	requestHeader() HttpHeader
	requestBody() HttpRequestBody
}

type HttpClient interface {
	Run(request HttpRequest) (HttpTest, error)
}

type HttpResponseStatus int

type ResponseTime time.Duration

type HttpTest interface {
	ExpectStatus(status HttpResponseStatus) TestResult
	ExpectTime(responseTime ResponseTime) TestResult
	ExpectHeader(name HttpHeaderName, value HttpHeaderValue) TestResult
	ExpectBodyContainsString(part string) TestResult
	ExpectBodyMatches(pattern string) TestResult
	ExpectBodySatisfies(predicate func(body string) bool) TestResult
}

type TestResult interface {
	Success() bool
	Comparison() Comparison
}

type Comparison interface {
	fmt.Stringer
	Expected() string
	Actual() string
}
