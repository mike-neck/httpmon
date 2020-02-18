package httpmon

import (
	"fmt"
	"io"
	"time"
)

type HttpRequestMethod string

const (
	GET    HttpRequestMethod = "GET"
	POST   HttpRequestMethod = "POST"
	PUT    HttpRequestMethod = "PUT"
	DELETE HttpRequestMethod = "DELETE"
)

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
	requestMethod() HttpRequestMethod
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
