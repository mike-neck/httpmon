package httpmon

import "fmt"

// RunTest runs http access test against given url, inspecting its response.
func RunTest(method, URL, timout string, status int) ([]TestResult, error) {
	caseRunner, err := NewCaseRunner(method, URL, timout)
	if err != nil {
		return nil, err
	}
	httpTest, err := caseRunner.Run()
	if err != nil {
		return nil, err
	}
	results := make([]TestResult, 0)
	statusResult := httpTest.ExpectStatus(HttpStatus(status))
	results = append(results, statusResult)
	return results, nil
}

type CaseRunner interface {
	Run() (HttpTest, error)
}

type defaultCaseRunner struct {
	HttpMon
	HttpTestRequest
}

func (c *defaultCaseRunner) Run() (HttpTest, error) {
	httpTest, err := c.HttpMon.Run(c.HttpTestRequest)
	if err != nil {
		return nil, &HttpCommunicationError{
			CaseError{
				Message:  "net work error",
				Original: err,
			},
		}
	}
	return httpTest, nil
}

type CaseError struct {
	Message  string
	Original error
}

type UserInputError struct {
	CaseError
}

func (err *UserInputError) Error() string {
	return fmt.Sprintf("input error: %s, %v", err.Message, err.Original)
}

type HttpCommunicationError struct {
	CaseError
}

func (err *HttpCommunicationError) Error() string {
	return fmt.Sprintf("%s: %v", err.Message, err.Original)
}

func NewCaseRunner(method, URL, timout string) (CaseRunner, error) {
	httpMethod, err := ToHttpMethod(method)
	if err != nil {
		return nil, &UserInputError{
			CaseError{
				Message:  "invalid http method",
				Original: err,
			},
		}
	}
	httpMon, err := NewHttpMon(timout)
	if err != nil {
		return nil, &UserInputError{
			CaseError{
				Message:  "invalid timeout",
				Original: err,
			},
		}
	}
	request, err := NewRequest(*httpMethod, URL)
	if err != nil {
		return nil, &UserInputError{
			CaseError{
				Message:  "invalid url",
				Original: err,
			},
		}
	}
	return &defaultCaseRunner{
		HttpMon:         *httpMon,
		HttpTestRequest: request,
	}, nil
}
