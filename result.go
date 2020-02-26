package httpmon

import (
	"fmt"
	"github.com/hako/durafmt"
	"strings"
	"time"
)

// implementations of Comparison

type HttpStatusSuccess struct {
	UserExpected HttpResponseStatus
	Response     HttpResponseStatus
}

func (suc *HttpStatusSuccess) Success() bool {
	return true
}

func (suc *HttpStatusSuccess) Comparison() Comparison {
	return suc
}

func (suc *HttpStatusSuccess) Expected() string {
	return fmt.Sprintf("status = %d", suc.UserExpected)
}

func (suc *HttpStatusSuccess) Actual() string {
	return fmt.Sprintf("status = %d", suc.Response)
}

func (suc *HttpStatusSuccess) String() string {
	return "ok"
}

type HttpStatusFailure struct {
	UserExpected HttpResponseStatus
	Response     HttpResponseStatus
}

func (fail *HttpStatusFailure) Success() bool {
	return false
}

func (fail *HttpStatusFailure) Comparison() Comparison {
	return fail
}

func (fail *HttpStatusFailure) Expected() string {
	return fmt.Sprintf("status = %d", fail.UserExpected)
}

func (fail *HttpStatusFailure) Actual() string {
	return fmt.Sprintf("status = %d", fail.Response)
}

func (fail *HttpStatusFailure) String() string {
	expected := fail.Expected()
	actual := fail.Actual()
	return comparisonFailureString(expected, actual)
}

func comparisonFailureString(expected, actual string) string {
	return fmt.Sprintf(comparisonFailureTemplate, expected, actual)
}

var comparisonFailureTemplate string = `expected: %s
actual  : %s`

func newHttpStatusTestResult(userExpected, actualResponse HttpResponseStatus) TestResult {
	if userExpected == actualResponse {
		return &HttpStatusSuccess{
			UserExpected: userExpected,
			Response:     actualResponse,
		}
	} else {
		return &HttpStatusFailure{
			UserExpected: userExpected,
			Response:     actualResponse,
		}
	}
}

////////

type SoftHeaderTest struct {
	Name                HttpHeaderName
	ActualValues        HttpHeaderValues
	ExpectedHeaderValue HttpHeaderValue
}

func (h *SoftHeaderTest) Success() bool {
	return h.isSuccess()
}

func (h *SoftHeaderTest) Comparison() Comparison {
	return h
}

func (h *SoftHeaderTest) isSuccess() bool {
	exp := h.ExpectedHeaderValue
	for _, v := range h.ActualValues {
		if v == exp {
			return true
		}
	}
	return false
}

func (h *SoftHeaderTest) String() string {
	if h.isSuccess() {
		return "ok"
	} else {
		return comparisonFailureString(h.Expected(), h.Actual())
	}
}

func (h *SoftHeaderTest) Expected() string {
	return fmt.Sprintf("header[%s] value = '%s'", h.Name, h.ExpectedHeaderValue)
}

func (h *SoftHeaderTest) Actual() string {
	if len(h.ActualValues) == 0 {
		return fmt.Sprintf("header[%s] not found", h.Name)
	}
	values := make([]string, len(h.ActualValues))
	for i, v := range h.ActualValues {
		values[i] = fmt.Sprintf("'%s'", string(v))
	}
	str := strings.Join(values, ",")
	return fmt.Sprintf("header[%s] values = [%s]", h.Name, str)
}

////////

// ResponseTimeComparison

type ResponseTimeTest struct {
	ActualTime ResponseTime
	ExpectTime ResponseTime
}

func (res *ResponseTimeTest) Success() bool {
	return res.isSuccess()
}

func (res *ResponseTimeTest) Comparison() Comparison {
	return res
}

func (res *ResponseTimeTest) isSuccess() bool {
	return res.ActualTime <= res.ExpectTime
}

func (res *ResponseTimeTest) String() string {
	if res.isSuccess() {
		return "ok"
	}
	return comparisonFailureString(res.Expected(), res.Actual())
}

func (res *ResponseTimeTest) Expected() string {
	dt := durafmt.Parse(time.Duration(res.ExpectTime))
	return fmt.Sprintf("response = %s", dt.String())
}

func (res *ResponseTimeTest) Actual() string {
	dt := durafmt.Parse(time.Duration(res.ActualTime))
	return fmt.Sprintf("response = %s", dt.String())
}
