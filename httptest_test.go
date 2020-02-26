package httpmon

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestHttpTest_ExpectStatus_Success(t *testing.T) {
	headers := make(http.Header, 0)
	headers.Add("content-type", "application/json; charset=utf-8")
	var test HttpTest = &DefaultHttpTest{
		Status:       HttpResponseStatus(200),
		Header:       headers,
		ResponseTime: 1000,
	}

	result := test.ExpectStatus(200)

	assert.IsType(t, new(HttpStatusSuccess), result)
	assert.True(t, result.Success())

	result = test.ExpectHeader("Content-Type", "application/json; charset=utf-8")
	assert.NotNil(t, result)
	assert.True(t, result.Success())

	assert.Equal(t, ResponseTime(1000), test.Performance())

	within := test.ExpectResponseTimeWithin(1001)
	if within == nil {
		assert.Fail(t, "unexpected nil of test")
		return
	}
	assert.True(t, within.Success())
}

func TestHttpTest_ExpectStatus_Failure(t *testing.T) {
	var test HttpTest = &DefaultHttpTest{
		Status: HttpResponseStatus(200),
	}

	result := test.ExpectStatus(401)

	assert.IsType(t, new(HttpStatusFailure), result)
	assert.False(t, result.Success())
}

func TestHttpTest_ExpectHeader_NotFound(t *testing.T) {
	headers := make(http.Header, 0)
	test := &DefaultHttpTest{
		Status: 200,
		Header: headers,
	}

	result := test.ExpectHeader("content-type", "application/xml")
	assert.NotNil(t, result)
	assert.False(t, result.Success())
	comparison := result.Comparison()
	assert.Equal(t, "header[content-type] not found", comparison.Actual())
}
