package httpmon

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
