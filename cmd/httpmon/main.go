package main

import (
	"bytes"
	"fmt"
	"github.com/mike-neck/httpmon"
	"github.com/urfave/cli/v2"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ExitCode int

const (
	NoError        ExitCode = 0
	TestError      ExitCode = 1
	ExecutionError ExitCode = 2
	UserError      ExitCode = 4
)

var Options = struct {
	Method        string
	Timeout       string
	Status        string
	ResponseTime  string
	ExpectHeader  string
	RequestHeader string
}{
	Method:        "method",
	Timeout:       "timeout",
	Status:        "status",
	ResponseTime:  "response-time",
	ExpectHeader:  "expect-header",
	RequestHeader: "request-header",
}

func main() {
	app := createApplication()
	code := runApplication(app, os.Args)
	os.Exit(int(code))
}

func createApplication() *cli.App {
	var method string
	var timeout string
	var responseTime string
	var status int
	return &cli.App{
		Name:  "httpmon",
		Usage: "runs synthetic monitoring test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        Options.Method,
				Aliases:     []string{"X"},
				Usage:       "http method",
				Required:    false,
				Value:       "GET",
				Destination: &method,
			},
			&cli.StringFlag{
				Name:        Options.Timeout,
				Aliases:     []string{"t"},
				Usage:       "timeout - format: nU(n : positive int, U : unit s(seconds) or m(minutes))",
				Required:    false,
				Value:       "5s",
				Destination: &timeout,
			},
			&cli.IntFlag{
				Name:        Options.Status,
				Aliases:     []string{"s"},
				Usage:       "expected http status",
				Required:    false,
				Value:       200,
				Destination: &status,
			},
			&cli.StringFlag{
				Name:        Options.ResponseTime,
				Aliases:     []string{"r"},
				Usage:       "expected response time - format: nU(n : positive int, U : unit s(seconds) or m(minutes))",
				Required:    false,
				Value:       "5s",
				Destination: &responseTime,
			},
			&cli.StringSliceFlag{
				Name:     Options.ExpectHeader,
				Aliases:  []string{"eh"},
				Usage:    "expecting header (format name=value)",
				Required: false,
			},
			&cli.StringSliceFlag{
				Name:     Options.RequestHeader,
				Aliases:  []string{"H"},
				Usage:    "request header(format name=value)",
				Required: false,
			},
		},
		Action: func(context *cli.Context) error {
			result := createAction(method, timeout, responseTime, status)(context)
			result.Describe()
			return result.AsError()
		},
	}
}

type Action func(context *cli.Context) ActionResult

type ActionResult interface {
	AsError() error
	Describe()
}

type InvalidUserInputSuspension struct {
	Errors []error
}

func (sus *InvalidUserInputSuspension) AsError() error {
	return sus
}

func (sus *InvalidUserInputSuspension) Describe() {
	fmt.Println(sus.Error())
}

func (sus *InvalidUserInputSuspension) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("Input error\n")
	for i, err := range sus.Errors {
		buffer.WriteString(fmt.Sprintf("%2d %s", i+1, err.Error()))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

type ErrorInExecution struct {
	url      httpmon.HttpRequestURL
	method   string
	timeout  string
	delegate error
}

func (err *ErrorInExecution) AsError() error {
	return err
}

func (err *ErrorInExecution) Describe() {
	fmt.Println(err.Error())
}

func (err *ErrorInExecution) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("Execution error %s %s", err.method, err.url))
	buffer.WriteString("\n")
	if gerr, ok := err.delegate.(*httpmon.GoStandardError); ok {
		if gerr.IsTimeout() {
			buffer.WriteString(fmt.Sprintf("timeout: %s", err.timeout))
		} else {
			buffer.WriteString(gerr.Error())
		}
	} else {
		buffer.WriteString(err.Error())
	}
	return buffer.String()
}

type TestFailed struct {
	Test
	Failure []httpmon.Comparison
}

func (f *TestFailed) AsError() error {
	return f
}

func (f *TestFailed) Describe() {
	fmt.Print(f.Error())
}

func (f *TestFailed) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s failed (success: %d / tests: %d)", f.url, f.count-len(f.Failure), f.count))
	buffer.WriteString("\n")
	for i, c := range f.Failure {
		buffer.WriteString(fmt.Sprintf("%d:\n", i+1))
		buffer.WriteString(c.String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

type Test struct {
	url   httpmon.HttpRequestURL
	count int
}

func (s *Test) AsError() error {
	return nil
}

func (s *Test) Describe() {
	fmt.Printf("%s ok (success: %d / tests: %d) \n", s.url, s.count, s.count)
}

func createAction(method, timeout, response string, status int) Action {
	es := make([]error, 0)
	requestMethod, err := httpmon.NewHttpRequestMethod(method)
	if err != nil {
		es = append(es, err)
	}

	timeoutTime, err := parseTimeout(timeout)
	if err != nil {
		es = append(es, err)
	}
	config := httpmon.Config{
		RequestTimeout: timeoutTime,
	}

	responseTime, err := parseResponseTime(response)
	if err != nil {
		es = append(es, err)
	}

	if status < 100 || 600 <= status {
		es = append(es, &httpmon.UserError{
			ItemName:   "http status",
			Reason:     "invalid status code is present",
			InputValue: status,
		})
	}
	expectStatus := httpmon.ExpectStatusOf(status)

	return func(context *cli.Context) ActionResult {
		ehs := context.StringSlice(Options.ExpectHeader)
		expectedHeaders, err := parseExpectedHeaders(ehs)
		if err != nil {
			es = append(es, err)
		}

		rhs := context.StringSlice(Options.RequestHeader)
		requestHeaders, err := parseRequestHeaders(rhs)
		if err != nil {
			es = append(es, err)
		}

		if context.NArg() != 1 {
			es = append(es, &httpmon.UserError{
				ItemName:   "url",
				Reason:     "no url is given",
				InputValue: "<no url>",
			})
		}
		urlFactory := func() httpmon.HttpRequestURL {
			url := context.Args().First()
			return httpmon.HttpRequestURL(url)
		}

		if 0 < len(es) {
			return &InvalidUserInputSuspension{Errors: es}
		}

		url := urlFactory()
		testCase := httpmon.Case{
			ClientBuilder:        &config,
			HttpRequestMethod:    requestMethod,
			URL:                  url,
			RequestHeaders:       requestHeaders,
			ExpectStatus:         expectStatus,
			ExpectedHeaders:      expectedHeaders,
			ExpectedResponseTime: responseTime,
		}
		result, err := testCase.Run()
		if err != nil {
			return &ErrorInExecution{
				url:      url,
				method:   method,
				timeout:  timeout,
				delegate: err,
			}
		}

		count := result.TestCount
		test := Test{
			url:   url,
			count: count,
		}

		if result.Success {
			return &test
		}

		return &TestFailed{
			Test:    test,
			Failure: result.Failed,
		}
	}
}

var pattern *regexp.Regexp = regexp.MustCompile("^[0-9]+[sm]$")
var num *regexp.Regexp = regexp.MustCompile("^[0-9]+")
var unit *regexp.Regexp = regexp.MustCompile("[sm]$")

func parseTimeout(timeout string) (httpmon.Timeout, error) {
	if strings.TrimSpace(timeout) == "" {
		return httpmon.Timeout(3 * time.Second), nil
	}
	t, err := parseTime(Options.Timeout, timeout)
	if err != nil {
		return 0, err
	}
	if t.number == 0 {
		return httpmon.Timeout(3 * time.Second), nil
	}
	return httpmon.Timeout(t.ToDuration()), nil
}

func parseResponseTime(res string) (httpmon.ExpectedResponseTime, error) {
	if strings.TrimSpace(res) == "" {
		return httpmon.ExpectedResponseTimeOf(1 * time.Second), nil
	}
	r, err := parseTime(Options.ResponseTime, res)
	if err != nil {
		return httpmon.ExpectedResponseTimeOf(0), err
	}
	if r.number == 0 {
		return httpmon.ExpectedResponseTimeOf(3 * time.Second), nil
	}
	return httpmon.ExpectedResponseTimeOf(r.ToDuration()), nil
}

type Time struct {
	number int64
	unit   TimeUnit
}

func (t *Time) ToDuration() time.Duration {
	switch t.unit {
	case Sec:
		return time.Duration(t.number) * time.Second
	case Min:
		return time.Duration(t.number) * time.Minute
	}
	panic(fmt.Sprintf("invalid time: %v", t))
}

func (t *Time) String() string {
	return fmt.Sprintf("%d %s", t.number, t.unit)
}

type TimeUnit string

const (
	Sec TimeUnit = "s"
	Min TimeUnit = "m"
)

func TimeUnitFromString(s string) (TimeUnit, error) {
	switch TimeUnit(s) {
	case Sec:
		return Sec, nil
	case Min:
		return Min, nil
	}
	return "", fmt.Errorf("unknown time unit: %s", s)
}

func parseTime(itemName, t string) (*Time, error) {
	if !pattern.MatchString(t) {
		return nil, &httpmon.UserError{
			ItemName:   itemName,
			Reason:     "invalid format, expected numberUNIT format(ex. 20s number:20, unit:s, means 20 sec)",
			InputValue: t,
		}
	}
	n := num.FindString(t)
	if n == "" {
		return nil, &httpmon.UserError{
			ItemName:   itemName,
			Reason:     "invalid number, expected numberUNIT format(ex. 20s number:20, unit:s, means 20 sec)",
			InputValue: t,
		}
	}
	number, err := strconv.ParseInt(n, 10, 0)
	if err != nil || number == 0 {
		return nil, &httpmon.UserError{
			ItemName:   itemName,
			Reason:     "invalid number, expected numberUNIT format(ex. 20s number:20, unit:s, means 20 sec)",
			InputValue: t,
		}
	}
	u := unit.FindString(t)
	if u == "" {
		return nil, &httpmon.UserError{
			ItemName:   itemName,
			Reason:     "invalid unit, expected numberUNIT format(ex. 20s number:20, unit:s, means 20 sec)",
			InputValue: t,
		}
	}

	timeUnit, err := TimeUnitFromString(u)
	if err != nil {
		return nil, &httpmon.UserError{
			ItemName:   itemName,
			Reason:     "invalid unit, expected numberUNIT format(ex. 20s number:20, unit:s, means 20 sec)",
			InputValue: t,
		}
	}

	return &Time{
		number: number,
		unit:   timeUnit,
	}, nil
}

func parseExpectedHeaders(headers []string) ([]httpmon.ExpectedHeader, error) {
	hs, err := parseHeaders(Options.ExpectHeader, headers)
	if err != nil {
		return nil, err
	}
	ehs := make([]httpmon.ExpectedHeader, len(hs))
	for i, h := range hs {
		ehs[i] = httpmon.ExpectedHeader{
			Name:  h.Name,
			Value: h.Value,
		}
	}
	return ehs, nil
}

func parseRequestHeaders(headers []string) ([]httpmon.RequestHeader, error) {
	hs, err := parseHeaders(Options.RequestHeader, headers)
	if err != nil {
		return nil, err
	}
	rhs := make([]httpmon.RequestHeader, len(hs))
	for i, h := range hs {
		rhs[i] = httpmon.RequestHeader{
			Name:  h.Name,
			Value: h.Value,
		}
	}
	return rhs, nil
}

type Header struct {
	Name  httpmon.HttpHeaderName
	Value httpmon.HttpHeaderValue
}

func parseHeaders(itemName string, headers []string) ([]Header, error) {
	hs := make([]Header, 0)
	for _, h := range headers {
		s := strings.Split(h, "=")
		if len(s) > 2 {
			return nil, &httpmon.UserError{
				ItemName:   itemName,
				Reason:     "invalid format, expected header-name=header-value(ex. accept=application/json)",
				InputValue: h,
			}
		}
		n := strings.TrimSpace(s[0])
		if n == "" {
			return nil, &httpmon.UserError{
				ItemName:   itemName,
				Reason:     "invalid header name, expected header-name=header-value(ex. accept=application/json)",
				InputValue: h,
			}
		}
		header := Header{
			Name:  httpmon.HttpHeaderName(n),
			Value: httpmon.HttpHeaderValue(strings.TrimSpace(s[1])),
		}
		hs = append(hs, header)
	}
	return hs, nil
}

func runApplication(app *cli.App, arguments []string) ExitCode {
	err := app.Run(arguments)
	if err != nil {
		if _, ok := err.(*InvalidUserInputSuspension); ok {
			return UserError
		}
		if _, ok := err.(*ErrorInExecution); ok {
			return ExecutionError
		}
		if _, ok := err.(*TestFailed); ok {
			return TestError
		}
		return ExecutionError
	}
	return NoError
}
