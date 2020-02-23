package httpmon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpTest_ExpectStatus_Success(t *testing.T) {
	var test HttpTest = &defaultHttpTest{
		Status: HttpResponseStatus(200),
	}

	result := test.ExpectStatus(200)

	assert.IsType(t, new(HttpStatusSuccess), result)
	assert.True(t, result.Success())
}

func TestHttpTest_ExpectStatus_Failure(t *testing.T) {
	var test HttpTest = &defaultHttpTest{
		Status: HttpResponseStatus(200),
	}

	result := test.ExpectStatus(401)

	assert.IsType(t, new(HttpStatusFailure), result)
	assert.False(t, result.Success())
}
