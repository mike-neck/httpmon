package main

import (
	"errors"
	"fmt"
	"github.com/mike-neck/httpmon"
	"github.com/urfave/cli/v2"
	"os"
)

type ExitCode int

const (
	NoError        ExitCode = 0
	TestError      ExitCode = 1
	ExecutionError ExitCode = 2
)

func main() {
	app := createApplication(testCaseAction)
	code := runApplication(app, os.Args...)
	os.Exit(int(code))
}

type ActionFactory func(method, timeout *string, status *int) cli.ActionFunc

func createApplication(actionFactory ActionFactory) *cli.App {
	var method string
	var timeout string
	var status int
	return &cli.App{
		Name:  "httpmon",
		Usage: "runs synthetic monitoring test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "method",
				Aliases:     []string{"X"},
				Usage:       "http method",
				Required:    false,
				Value:       "GET",
				Destination: &method,
			},
			&cli.StringFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Usage:       "timeout - format: nU(n : positive int, U : unit s(seconds) or m(minutes))",
				Required:    false,
				Value:       "5s",
				Destination: &timeout,
			},
			&cli.IntFlag{
				Name:        "status",
				Aliases:     []string{"s"},
				Usage:       "expected http status",
				Required:    false,
				Value:       200,
				Destination: &status,
			},
		},
		Action: actionFactory(&method, &timeout, &status),
	}
}

func runApplication(app *cli.App, arguments ...string) ExitCode {
	err := app.Run(arguments)
	if err != nil {
		if e, ok := err.(*TestFailed); ok {
			fmt.Println(e.Error())
			return TestError
		} else if e, ok := err.(*httpmon.UserInputError); ok {
			fmt.Println(e.Error())
			return ExecutionError
		} else if e, ok := err.(*httpmon.HttpCommunicationError); ok {
			fmt.Println(e.Error())
			return ExecutionError
		} else {
			fmt.Println(err)
			return ExecutionError
		}
	}
	return NoError
}

type TestFailed struct {
	tests  int
	failed int
}

func (f *TestFailed) Error() string {
	return fmt.Sprintf("test failed: %d failed in %d cases", f.failed, f.tests)
}

func testCaseAction(method, timeout *string, status *int) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if ctx.NArg() != 1 {
			return &httpmon.UserInputError{
				CaseError: httpmon.CaseError{
					Message:  "url required",
					Original: errors.New("missing url"),
				},
			}
		}
		url := ctx.Args().First()
		results, err := httpmon.RunGetRequestCase(*method, url, *timeout, *status)
		if err != nil {
			return err
		}
		failed := 0
		for _, result := range results {
			if !result.IsSuccess() {
				comparison := result.Comparison()
				fmt.Println(comparison.ItemName)
				fmt.Printf("  expect: %v\n", comparison.Expect)
				fmt.Printf("  actual: %v\n", comparison.Actual)
				failed++
			}
		}
		if failed > 0 {
			return &TestFailed{
				tests:  len(results),
				failed: failed,
			}
		}
		fmt.Println("ok")
		return nil
	}
}
