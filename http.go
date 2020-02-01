package httpmon

import (
	"fmt"
	"strings"
)

type HttpMon struct {
	TimeOut
	HttpClientFactory
}

func NewHttpMon(timeout string) (*HttpMon, error) {
	timeOut, err := TimeOutFromString(timeout)
	if err != nil {
		return nil, err
	}
	factory := func(out TimeOut) HttpClient {
		return DefaultHttpClient(out)
	}
	return &HttpMon{
		TimeOut:           *timeOut,
		HttpClientFactory: factory,
	}, nil
}

type HttpMethod string

const (
	GET    HttpMethod = "GET"
	POST   HttpMethod = "POST"
	PUT    HttpMethod = "PUT"
	DELETE HttpMethod = "DELETE"
)

func ToHttpMethod(input string) (*HttpMethod, error) {
	method := HttpMethod(strings.ToUpper(input))
	switch method {
	case GET:
		return &method, nil
	case POST:
	case PUT:
	case DELETE:
		return nil, fmt.Errorf("%s is not supported", method)
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}

type HttpHeaders map[string][]string

type Queries map[string][]string

type HttpTestRequest interface {
	Method() HttpMethod
	URL() string
}

type TestResult interface {
	IsSuccess() bool
	Comparison() Comparison
}

type Comparison struct {
	ItemName string
	Expect   interface{}
	Actual   interface{}
}

type HttpStatus int

type HttpTest interface {
	ExpectStatus(status HttpStatus) TestResult
}

func (hm *HttpMon) httpClient() HttpClient {
	return hm.HttpClientFactory(hm.TimeOut)
}

func (hm *HttpMon) Run(request HttpTestRequest) (HttpTest, error) {
	client := hm.httpClient()
	response, err := client.Run(request.Method(), request.URL())
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Close() }()
	return &test{status: response.StatusCode()}, nil
}

type test struct {
	status HttpStatus
}

func (t *test) ExpectStatus(status HttpStatus) TestResult {
	return &httpStatusComparison{
		result: t.status == status,
		expect: status,
		actual: t.status,
	}
}

type httpStatusComparison struct {
	result bool
	expect HttpStatus
	actual HttpStatus
}

func (h *httpStatusComparison) IsSuccess() bool {
	return h.result
}

func (h *httpStatusComparison) Comparison() Comparison {
	return Comparison{
		ItemName: "status",
		Expect:   h.expect,
		Actual:   h.actual,
	}
}
