package httpmon

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpRequestMethod(t *testing.T) {
	assert.Equal(t, HttpRequestMethod("GET"), GET)
	assert.Equal(t, HttpRequestMethod("POST"), POST)
}
