package httpmon

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRunTestCase_InvalidInput(t *testing.T) {
	results, err := RunTestCase("HEAD", "https://examaple.com", "5s", 404)
	assert.NotNil(t, err)
	assert.IsType(t, new(UserInputError), err)
	assert.Nil(t, results)
}

func TestRunTestCase_HttpCommunicationError(t *testing.T) {
	results, err := RunTestCase("GET", "http://localhost:4000/test", "2s", 200)
	assert.NotNil(t, err)
	assert.IsType(t, new(HttpCommunicationError), err)
	assert.Nil(t, results)
}

func TestRunTestCase(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = fmt.Fprintln(w, "ok")
	}))
	defer func() { _ = server.Close }()

	results, err := RunTestCase("GET", server.URL, "2s", 200)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(results))
	assert.True(t, results[0].IsSuccess())
}
