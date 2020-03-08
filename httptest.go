package httpmon

import (
	"net/http"
	"net/textproto"
)

type DefaultHttpTest struct {
	Status HttpResponseStatus
	Header http.Header
	ResponseTime
}

func (dht *DefaultHttpTest) Performance() ResponseTime {
	return dht.ResponseTime
}

func (dht *DefaultHttpTest) ExpectResponseTimeWithin(responseTime ResponseTime) TestResult {
	return &ResponseTimeTest{
		ActualTime: dht.ResponseTime,
		ExpectTime: responseTime,
	}
}

func (dht *DefaultHttpTest) ExpectStatus(status HttpResponseStatus) TestResult {
	return newHttpStatusTestResult(status, dht.Status)
}

func (dht *DefaultHttpTest) ExpectTime(responseTime ResponseTime) TestResult {
	panic("implement me")
}

func (dht *DefaultHttpTest) ExpectHeader(name HttpHeaderName, value HttpHeaderValue) TestResult {
	s := string(name)
	key := textproto.CanonicalMIMEHeaderKey(s)

	header := textproto.MIMEHeader(dht.Header)
	values, ok := header[key]
	if !ok {
		return &SoftHeaderTest{
			Name:                name,
			ExpectedHeaderValue: value,
		}
	}

	hvs := make(HttpHeaderValues, len(values))
	for i, v := range values {
		hvs[i] = HttpHeaderValue(v)
	}
	return &SoftHeaderTest{
		Name:                name,
		ActualValues:        hvs,
		ExpectedHeaderValue: value,
	}
}

func (dht *DefaultHttpTest) ExpectBodyContainsString(part string) TestResult {
	panic("implement me")
}

func (dht *DefaultHttpTest) ExpectBodyMatches(pattern string) TestResult {
	panic("implement me")
}

func (dht *DefaultHttpTest) ExpectBodySatisfies(predicate func(body string) bool) TestResult {
	panic("implement me")
}
