package httpmon

import "fmt"

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
	template := `expected: %s
actual  : %s`
	return fmt.Sprintf(template, fail.Expected(), fail.Actual())
}
