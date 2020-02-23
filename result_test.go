package httpmon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComparison_HttpStatus_Success(t *testing.T) {
	var success Comparison
	success = &HttpStatusSuccess{
		UserExpected: 200,
		Response:     200,
	}

	assert.Equal(t, "ok", success.String())
	assert.Equal(t, "status = 200", success.Expected())
	assert.Equal(t, "status = 200", success.Actual())
}

func TestComparison_HttpStatus_Failure(t *testing.T) {
	var failure Comparison
	failure = &HttpStatusFailure{
		UserExpected: 200,
		Response:     401,
	}

	assert.Equal(t,
		`expected: status = 200
actual  : status = 401`, failure.String())
	assert.Equal(t, "status = 200", failure.Expected())
	assert.Equal(t, "status = 401", failure.Actual())
}
