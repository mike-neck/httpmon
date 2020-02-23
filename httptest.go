package httpmon

type defaultHttpTest struct {
	Status HttpResponseStatus
}

func (dht *defaultHttpTest) ExpectStatus(status HttpResponseStatus) TestResult {
	return newHttpStatusTestResult(status, dht.Status)
}

func (dht *defaultHttpTest) ExpectTime(responseTime ResponseTime) TestResult {
	panic("implement me")
}

func (dht *defaultHttpTest) ExpectHeader(name HttpHeaderName, value HttpHeaderValue) TestResult {
	panic("implement me")
}

func (dht *defaultHttpTest) ExpectBodyContainsString(part string) TestResult {
	panic("implement me")
}

func (dht *defaultHttpTest) ExpectBodyMatches(pattern string) TestResult {
	panic("implement me")
}

func (dht *defaultHttpTest) ExpectBodySatisfies(predicate func(body string) bool) TestResult {
	panic("implement me")
}
