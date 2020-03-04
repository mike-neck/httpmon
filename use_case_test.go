package httpmon

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestConfigNewClient(t *testing.T) {
	config := Config{RequestTimeout: 10000}

	client := config.newClient()

	assert.IsType(t, new(DefaultHttpClient), client)
}

func TestGetCase_NewRequest_WithoutHeader(t *testing.T) {
	config := &Config{RequestTimeout: Timeout(3 * time.Second)}
	c := Case{
		ClientBuilder:     config,
		HttpRequestMethod: GET,
		URL:               "https://example.com/test",
		RequestHeaders:    []RequestHeader{},
		ExpectStatus:      ExpectStatusOf(200),
		ExpectedHeaders:   []ExpectedHeader{},
	}

	request := c.newRequest()

	assert.Equal(t, GetMethod, request.requestMethod())
	assert.Equal(t, HttpRequestURL("https://example.com/test"), request.requestURL())
	assert.Equal(t, HttpHeader{}, request.requestHeader())
}

func TestGetCase_NewRequest_WithHeader(t *testing.T) {
	config := &Config{RequestTimeout: Timeout(3 * time.Second)}
	getCase := Case{
		ClientBuilder:     config,
		HttpRequestMethod: GET,
		URL:               "https://example.com/test",
		RequestHeaders: []RequestHeader{
			{
				Name:  "Accept",
				Value: "application/json",
			},
			{
				Name:  "authorization",
				Value: "Bearer 11aa22bb33cc44dd55ee6f",
			},
		},
		ExpectStatus:    ExpectStatusOf(200),
		ExpectedHeaders: []ExpectedHeader{},
	}

	request := getCase.newRequest()

	assert.Equal(t, GetMethod, request.requestMethod())
	assert.Equal(t, HttpRequestURL("https://example.com/test"), request.requestURL())
	assert.Equal(t, HttpHeader{
		"Accept": HttpHeaderValues{
			"application/json",
		},
		"authorization": HttpHeaderValues{
			"Bearer 11aa22bb33cc44dd55ee6f",
		},
	}, request.requestHeader())
}

func TestGetCase_Run_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	test := NewMockHttpTest(ctrl)
	test.EXPECT().ExpectStatus(gomock.Eq(HttpResponseStatus(200))).Return(&HttpStatusSuccess{
		UserExpected: 200,
		Response:     200,
	})
	test.EXPECT().
		ExpectHeader(HttpHeaderName("content-type"), HttpHeaderValue("application/json")).
		Return(&SoftHeaderTest{
			Name: "content-type",
			ActualValues: HttpHeaderValues{
				"application/json",
			},
			ExpectedHeaderValue: "application/json",
		})
	test.EXPECT().
		ExpectResponseTimeWithin(gomock.Any()).
		DoAndReturn(func(rt ResponseTime) TestResult {
			return &ResponseTimeTest{
				ActualTime: ResponseTime(2100 * time.Millisecond),
				ExpectTime: rt,
			}
		})

	httpClient := NewMockHttpClient(ctrl)
	httpClient.EXPECT().
		Run(gomock.Any()).Return(test, nil)

	builder := NewMockClientBuilder(ctrl)
	builder.EXPECT().
		newClient().Return(httpClient)

	getCase := Case{
		ClientBuilder:     builder,
		HttpRequestMethod: GET,
		URL:               "https://example.com",
		RequestHeaders:    []RequestHeader{},
		ExpectStatus:      ExpectStatusOf(200),
		ExpectedHeaders: []ExpectedHeader{
			{
				Name:  "content-type",
				Value: "application/json",
			},
		},
		ExpectedResponseTime: ExpectedResponseTimeOf(3 * time.Second),
	}

	caseResult, err := getCase.Run()
	if err != nil {
		assert.Fail(t, "unexpected error: %v", err)
		return
	}

	assert.True(t, caseResult.Success)
	assert.Len(t, caseResult.Failed, 0)
	assert.Equal(t, 3, caseResult.TestCount)
}

func TestGetCase_Run_Error(t *testing.T) {
	expectedError := fmt.Errorf("http error")

	ctrl := gomock.NewController(t)
	httpClient := NewMockHttpClient(ctrl)
	httpClient.EXPECT().
		Run(gomock.Any()).Return(nil, expectedError)
	builder := NewMockClientBuilder(ctrl)
	builder.EXPECT().
		newClient().Return(httpClient)

	getCase := Case{
		ClientBuilder:        builder,
		HttpRequestMethod:    GET,
		URL:                  "https://example.com",
		RequestHeaders:       []RequestHeader{},
		ExpectStatus:         ExpectStatusOf(200),
		ExpectedHeaders:      []ExpectedHeader{},
		ExpectedResponseTime: ExpectedResponseTimeOf(3 * time.Second),
	}

	caseResult, err := getCase.Run()
	if err == nil {
		assert.Fail(t, "unexpected success")
		return
	}

	assert.Equal(t, expectedError, err)
	assert.Equal(t, 0, caseResult.TestCount)
	assert.Len(t, caseResult.Failed, 0)
	assert.False(t, caseResult.Success)
}

func TestCase_Run_Failure(t *testing.T) {

	ctrl := gomock.NewController(t)
	httpClient := NewMockHttpClient(ctrl)
	httpClient.EXPECT().
		Run(gomock.Any()).Return(&DefaultHttpTest{
		Status: 404,
		Header: http.Header{
			"content-type": []string{"application/json"},
		},
		ResponseTime: ResponseTime(3 * time.Second),
	}, nil)
	builder := NewMockClientBuilder(ctrl)
	builder.EXPECT().
		newClient().Return(httpClient)

	getCase := Case{
		ClientBuilder:     builder,
		HttpRequestMethod: GET,
		URL:               "https://example.com",
		RequestHeaders:    []RequestHeader{},
		ExpectStatus:      ExpectStatusOf(200),
		ExpectedHeaders: []ExpectedHeader{
			{
				Name:  "Content-Type",
				Value: "application/xml",
			},
		},
	}

	result, err := getCase.Run()

	if err != nil {
		assert.Fail(t, "unexpected err: %v", err)
		return
	}
	assert.False(t, result.Success)
	assert.Len(t, result.Failed, 2)
}
