package main

import (
	"errors"
	"github.com/mike-neck/httpmon"
	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli/v2"
	"testing"
)

func TestRunApplication_TestFailed(t *testing.T) {
	app := &cli.App{
		Name: "testing",
		Action: func(c *cli.Context) error {
			return &TestFailed{
				tests:  5,
				failed: 3,
			}
		},
	}
	exitCode := runApplication(app, "test-app")
	assert.Equal(t, TestError, exitCode)
}

func TestRunApplication_UserInputError(t *testing.T) {
	app := &cli.App{
		Name: "testing",
		Action: func(c *cli.Context) error {
			return &httpmon.UserInputError{
				CaseError: httpmon.CaseError{
					Message:  "test",
					Original: errors.New("user input"),
				},
			}
		},
	}
	exitCode := runApplication(app, "test-app")
	assert.Equal(t, ExecutionError, exitCode)
}

func TestRunApplication_HttpCommunicationError(t *testing.T) {
	app := &cli.App{
		Name: "testing",
		Action: func(c *cli.Context) error {
			return &httpmon.HttpCommunicationError{
				CaseError: httpmon.CaseError{
					Message:  "test",
					Original: errors.New("user input"),
				},
			}
		},
	}
	exitCode := runApplication(app, "test-app")
	assert.Equal(t, ExecutionError, exitCode)
}

func TestRunApplication_UnknownError(t *testing.T) {
	app := &cli.App{
		Name: "testing",
		Action: func(c *cli.Context) error {
			return errors.New("unknown error")
		},
	}
	exitCode := runApplication(app, "test-app")
	assert.Equal(t, ExecutionError, exitCode)
}

func TestRunApplication_Success(t *testing.T) {
	app := &cli.App{
		Name: "testing",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
	exitCode := runApplication(app, "test-app")
	assert.Equal(t, NoError, exitCode)
}
