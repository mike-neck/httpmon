package httpmon

type CaseResult struct {
	Success   bool
	Failed    []Comparison
	TestCount int
}

func newEmptyCaseResult() CaseResult {
	return CaseResult{
		Success:   true,
		Failed:    []Comparison{},
		TestCount: 0,
	}
}

func (result *CaseResult) Append(testResult TestResult) {
	result.TestCount += 1
	if !testResult.Success() {
		result.Success = false
		result.Failed = append(result.Failed, testResult.Comparison())
	}
}

type ClientBuilder interface {
	newClient() HttpClient
}

type Config struct {
	RequestTimeout Timeout
}

func (c *Config) newClient() HttpClient {
	return NewHttpClient(c.RequestTimeout)
}

type RequestHeader struct {
	Name  HttpHeaderName
	Value HttpHeaderValue
}

type ExpectedHeader struct {
	Name  HttpHeaderName
	Value HttpHeaderValue
}

type Case struct {
	ClientBuilder
	HttpRequestMethod
	URL             HttpRequestURL
	RequestHeaders  []RequestHeader
	ExpectStatus    HttpResponseStatus
	ExpectedHeaders []ExpectedHeader
}

func (c *Case) Run() (CaseResult, error) {
	client := c.newClient()
	request := c.newRequest()
	test, err := client.Run(request)
	if err != nil {
		return CaseResult{}, err
	}
	result := newEmptyCaseResult()
	if c.ExpectStatus.IsValidValue() {
		status := test.ExpectStatus(c.ExpectStatus)
		result.Append(status)
	}
	for _, hdr := range c.ExpectedHeaders {
		headerResult := test.ExpectHeader(hdr.Name, hdr.Value)
		result.Append(headerResult)
	}
	return result, nil
}

func (c *Case) newRequest() HttpRequest {
	request := c.HttpRequestMethod(c.URL)
	for _, h := range c.RequestHeaders {
		request.AddHeader(h.Name, h.Value)
	}
	return request
}
