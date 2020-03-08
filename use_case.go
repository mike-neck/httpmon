package httpmon

import "time"

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

type ExpectStatus func() HttpResponseStatus

func ExpectStatusOf(status int) ExpectStatus {
	return func() HttpResponseStatus {
		return HttpResponseStatus(status)
	}
}

type RequestHeader struct {
	Name  HttpHeaderName
	Value HttpHeaderValue
}

type ExpectedHeader struct {
	Name  HttpHeaderName
	Value HttpHeaderValue
}

type ExpectedResponseTime func() ResponseTime

func ExpectedResponseTimeOf(t time.Duration) ExpectedResponseTime {
	return func() ResponseTime {
		return ResponseTime(t)
	}
}

type Case struct {
	ClientBuilder
	HttpRequestMethod
	URL            HttpRequestURL
	RequestHeaders []RequestHeader
	ExpectStatus
	ExpectedHeaders []ExpectedHeader
	ExpectedResponseTime
}

func (c *Case) Run() (CaseResult, error) {
	client := c.newClient()
	request := c.newRequest()
	test, err := client.Run(request)
	if err != nil {
		return CaseResult{}, err
	}
	result := newEmptyCaseResult()
	if c.ExpectStatus != nil {
		status := test.ExpectStatus(c.ExpectStatus())
		result.Append(status)
	}
	for _, hdr := range c.ExpectedHeaders {
		headerResult := test.ExpectHeader(hdr.Name, hdr.Value)
		result.Append(headerResult)
	}
	if c.ExpectedResponseTime != nil {
		timeResult := test.ExpectResponseTimeWithin(c.ExpectedResponseTime())
		result.Append(timeResult)
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
