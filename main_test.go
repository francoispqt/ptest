package main

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const SuccessPackageName = "github.com/francoispqt/ptest/successtestdir"
const ErrorPackageName = "github.com/francoispqt/ptest/errortestdir"

// Interface
type MockTest struct {
	PackageName string
}

func (m MockTest) Package() string {
	return m.PackageName
}

func (m MockTest) Run(dirPath string, testsResult chan TestResult, args []string) {
	if m.Package() == "error" {
		testsResult <- TestResult{[]byte("test"), errors.New("test error"), m.PackageName}
		return
	}
	testsResult <- TestResult{[]byte("test"), nil, m.PackageName}
}

func TestNewTests(t *testing.T) {
	t.Run("It should panic if dir does not exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				require.True(t, strings.Contains(r.(error).Error(), "no such file or directory"))
			}
		}()
		testsResultChan := make(chan TestResult)
		_ = NewTests("somedir", testsResultChan, []string{})
		// should not be called
		require.Equal(t, 0, 1)
	})
	t.Run("It should return the list of tests", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(SuccessPackageName, testsResultChan, []string{})
		require.Equal(t, len(tests), 3)
	})
}

func TestGetTestsResult(t *testing.T) {
	t.Run("It should return the right exit code", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(SuccessPackageName, testsResultChan, []string{})
		// check results
		success := GetTestsResult(testsResultChan, tests)
		require.Equal(t, success, 0)
	})

	t.Run("It should return the right exit code", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(ErrorPackageName, testsResultChan, []string{})
		// check results
		success := GetTestsResult(testsResultChan, tests)
		require.Equal(t, success, 1)
	})
}

func TestGetArgs(t *testing.T) {
	t.Run("It should return the second arg and the rest of args", func(t *testing.T) {
		f, r := getArgs([]string{"first", "second", "rest1", "rest2"})
		require.Equal(t, f, "second")
		require.Equal(t, r[0], "rest1")
		require.Equal(t, r[1], "rest2")
	})
}
