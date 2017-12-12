package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const SuccessPackageName = "github.com/francoispqt/ptest/successtestdir"
const ErrorPackageName = "github.com/francoispqt/ptest/errortestdir"
const SkipPackageName = "github.com/francoispqt/ptest/skiptestdir"

func TestNewTests(t *testing.T) {
	t.Run("It should panic if dir does not exist", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				require.True(t, strings.Contains(r.(error).Error(), "no such file or directory"))
			}
		}()
		testsResultChan := make(chan TestResult)
		_ = NewTests("somedir", testsResultChan, []string{}, false)
		// should not be called
		require.Equal(t, 0, 1)
	})
	t.Run("It should return the list of tests", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(SuccessPackageName, testsResultChan, []string{}, false)
		require.Equal(t, len(tests), 3)
	})

	t.Run("It should skip the tests because there is a .skiptest", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(SkipPackageName, testsResultChan, []string{}, true)
		require.Equal(t, len(tests), 0)
	})
}

func TestGetTestsResult(t *testing.T) {
	t.Run("It should return the right exit code", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(SuccessPackageName, testsResultChan, []string{}, false)
		// check results
		success := GetTestsResult(testsResultChan, tests)
		require.Equal(t, success, 0)
	})

	t.Run("It should return the right exit code", func(t *testing.T) {
		testsResultChan := make(chan TestResult)
		tests := NewTests(ErrorPackageName, testsResultChan, []string{}, false)
		// check results
		success := GetTestsResult(testsResultChan, tests)
		require.Equal(t, success, 1)
	})
}

func TestGetArgs(t *testing.T) {
	t.Run("It should return the second arg and the rest of args", func(t *testing.T) {
		f, r, err := getArgs([]string{"first", "second", "rest1", "rest2"})
		require.Nil(t, err)
		require.Equal(t, f, "second")
		require.Equal(t, r[0], "rest1")
		require.Equal(t, r[1], "rest2")
	})

	t.Run("It should return the second arg and the rest of args", func(t *testing.T) {
		_, _, err := getArgs([]string{"first"})
		require.Error(t, err)
	})
}

func TestRun(t *testing.T) {
	t.Run("It should run the test suite and return the right signal", func(t *testing.T) {
		signal := Run(SuccessPackageName, []string{})
		require.Equal(t, signal, 0)
	})

	t.Run("It should run the test suite and return the right signal", func(t *testing.T) {
		signal := Run(ErrorPackageName, []string{})
		require.Equal(t, signal, 1)
	})

	t.Run("It should return 1 because no test are found", func(t *testing.T) {
		signal := Run(SkipPackageName, []string{})
		require.Equal(t, signal, 1)
	})
}
