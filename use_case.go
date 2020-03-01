package httpmon

type Case interface {
	Run() (CaseResult, error)
}

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

type GetCase struct {
	ClientBuilder
	URL             HttpRequestURL
	RequestHeaders  []RequestHeader
	ExpectStatus    HttpResponseStatus
	ExpectedHeaders []ExpectedHeader
}

func (getCase *GetCase) Run() (CaseResult, error) {
	client := getCase.newClient()
	request := getCase.newRequest()
	test, err := client.Run(request)
	if err != nil {
		return CaseResult{}, err
	}
	result := newEmptyCaseResult()
	if getCase.ExpectStatus.IsValidValue() {
		status := test.ExpectStatus(getCase.ExpectStatus)
		result.Append(status)
	}
	return result, nil
}

func (getCase *GetCase) newRequest() HttpRequest {
	request := GET(getCase.URL)
	for _, h := range getCase.RequestHeaders {
		request.AddHeader(h.Name, h.Value)
	}
	return request
}
