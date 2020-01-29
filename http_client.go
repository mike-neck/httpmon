package httpmon

type HttpMon struct {
}

type HttpMethod string

type HttpHeaders map[string][]string

type Queries map[string][]string

type HttpTestRequest interface {
	Method() HttpMethod
	Headers() HttpHeaders
	URL() string
	Queries() Queries
}

type TestResult int

const (
	TestSuccess TestResult = iota
	TestFailure
)

type HttpStatus int

type HttpHeaderName string

type HttpHeaderValue string

type HttpBodyTester func(body []byte) (result bool, err error)

type HttpTest interface {
	ExpectStatus(status HttpStatus) TestResult
	ExpectHeader(name HttpHeaderName) HttpHeaderTest
	ExpectBody(tester HttpBodyTester) TestResult
}

type HttpHeaderTest interface {
	Exists() TestResult
	HasValue(value HttpHeaderValue) TestResult
}
