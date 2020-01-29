package httpmon

type HttpMon struct {
}

type HttpMethod string

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
