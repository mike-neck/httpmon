package httpmon

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_HttpStatus_Success(t *testing.T) {
	success := &HttpStatusSuccess{
		UserExpected: 200,
		Response:     200,
	}

	assert.Equal(t, "ok", success.String())
	assert.Equal(t, "status = 200", success.Expected())
	assert.Equal(t, "status = 200", success.Actual())

	var result TestResult
	result = success
	assert.True(t, result.Success())
	assert.Equal(t, success, result.Comparison())
}

func Test_HttpStatus_Failure(t *testing.T) {
	failure := &HttpStatusFailure{
		UserExpected: 200,
		Response:     401,
	}

	assert.Equal(t,
		`expected: status = 200
actual  : status = 401`, failure.String())
	assert.Equal(t, "status = 200", failure.Expected())
	assert.Equal(t, "status = 401", failure.Actual())

	var result TestResult
	result = failure
	assert.False(t, result.Success())
	assert.Equal(t, failure, result.Comparison())
}

func Test_SoftHeaderTest_Success(t *testing.T) {
	var comparison Comparison
	comparison = &SoftHeaderTest{
		Name: HttpHeaderName("Content-Type"),
		ActualValues: HttpHeaderValues{
			"application/json",
			"application/vnd.v2.example.com+json",
		},
		ExpectedHeaderValue: HttpHeaderValue("application/json"),
	}

	assert.Equal(t, "header[Content-Type] values = ['application/json','application/vnd.v2.example.com+json']", comparison.Actual())
	assert.Equal(t, "header[Content-Type] value = 'application/json'", comparison.Expected())
	assert.Equal(t, "ok", comparison.String())
}

func Test_SoftHeaderTest_Failure(t *testing.T) {
	var comparison Comparison
	comparison = &SoftHeaderTest{
		Name: HttpHeaderName("Content-Type"),
		ActualValues: HttpHeaderValues{
			"application/xml",
		},
		ExpectedHeaderValue: HttpHeaderValue("application/json"),
	}

	assert.Equal(t, "header[Content-Type] values = ['application/xml']", comparison.Actual())
	assert.Equal(t, "header[Content-Type] value = 'application/json'", comparison.Expected())
	assert.Equal(t, `expected: header[Content-Type] value = 'application/json'
actual  : header[Content-Type] values = ['application/xml']`, comparison.String())
}

func Test_SoftHeaderTest_NotFound(t *testing.T) {
	var comparison Comparison
	comparison = &SoftHeaderTest{
		Name:                HttpHeaderName("Content-Type"),
		ExpectedHeaderValue: HttpHeaderValue("application/json"),
	}

	assert.Equal(t, "header[Content-Type] not found", comparison.Actual())
	assert.Equal(t, "header[Content-Type] value = 'application/json'", comparison.Expected())
	assert.Equal(t, `expected: header[Content-Type] value = 'application/json'
actual  : header[Content-Type] not found`, comparison.String())
}

func TestResponseTimeTest_Success(t *testing.T) {
	var comparison Comparison
	comparison = &ResponseTimeTest{
		ActualTime: ResponseTime(3280 * time.Millisecond),
		ExpectTime: ResponseTime(5 * time.Second),
	}

	assert.Equal(t, "response = 3 seconds 280 milliseconds", comparison.Actual())
	assert.Equal(t, "response = 5 seconds", comparison.Expected())
	assert.Equal(t, "ok", comparison.String())
}
