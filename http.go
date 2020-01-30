package httpmon

type HttpMon struct {
	TimeOut
	HttpClientFactory
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

func (t test) ExpectStatus(status HttpStatus) TestResult {
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
