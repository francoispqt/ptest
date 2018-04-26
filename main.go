// Package ptest is small package used to run tests in package and subpackages concurrently
// adds a slightly better output
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
)

// constants
const (
	vendorDir    = "vendor"
	SkipFileName = ".ptestskip"
)

// variables, goPath, colors and al
var (
	goPath = os.Getenv("GOPATH")
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func getArgs(args []string) (string, []string, error) {
	if len(args) < 2 {
		return "", []string{}, errors.New("Usage ptest ${packageImportPath}")
	}
	return args[1], args[2:], nil
}

// TestRunner interface for testing
type TestRunner interface {
	Run(chan TestResult, []string)
	Package() string
}

// Test is the structure representing a package test
type Test struct {
	PackageName string
}

// Run runs test for the given directory and writes result to a channel
func (t Test) Run(testsResult chan TestResult, args []string) {
	cmdArgs := append([]string{"test", t.Package()}, args...)
	log.Printf("Testing %s", yellow(t.Package()))
	out, err := exec.Command("go", cmdArgs...).Output()
	testsResult <- TestResult{out, err, t.Package()}
}

// Package return the name of the package
func (t Test) Package() string {
	return t.PackageName
}

// TestResult is the structure holding the result of the call to go test
type TestResult struct {
	Out     []byte
	Err     error
	Package string
}

// NewTests gets a list of the testable packages and subpackages and runs tests
func NewTests(dirPath string, testsResultChan chan TestResult, args []string, allowSkip bool) []TestRunner {
	tests := make([]TestRunner, 0)
	files, err := ioutil.ReadDir(goPath + "/src/" + dirPath)
	if err != nil {
		panic(err)
	}
	hasTest := false
	skip := false
	for _, f := range files {
		// avoid testing
		if f.IsDir() && f.Name() != vendorDir {
			// should run concurrently
			tests = append(tests, NewTests(dirPath+"/"+f.Name(), testsResultChan, args, allowSkip)...)
		} else if strings.HasSuffix(f.Name(), "_test.go") && !hasTest && !skip {
			test := Test{dirPath}
			tests = append(tests, test)
			// run test
			go test.Run(testsResultChan, args)
			hasTest = true
		} else if f.Name() == SkipFileName && allowSkip {
			skip = true
		}
	}
	return tests
}

// GetTestsResult gets the result from tests, the result is the exit code
func GetTestsResult(testsResultChan chan TestResult, nTests []TestRunner) int {
	success := 0
	count := 0
	for testResult := range testsResultChan {
		if testResult.Err != nil {
			log.Printf(red("× %s"), testResult.Package)
			fmt.Println(string(testResult.Out))
			success = 1
		} else {
			log.Printf(green("✓ %s"), testResult.Package)
			fmt.Println(string(testResult.Out))
		}
		count++
		// close channel when all tests have been covered
		if count == len(nTests) {
			close(testsResultChan)
		}
	}
	// check result and exit
	if success == 0 {
		log.Print(green("Tests successfull !"))
	} else {
		log.Print(red("Tests fails !"))
	}
	return success
}

// Run runs the tests suite
func Run(packageRoot string, args []string) int {
	// create chan an defer its closure
	testsResultChan := make(chan TestResult)
	// create tests by recursively getting subpackages
	tests := NewTests(packageRoot, testsResultChan, args, true)
	if len(tests) == 0 {
		log.Print(red("No tests to run !"))
		return 1
	}
	// get tests results fron channel
	return GetTestsResult(testsResultChan, tests)
}

func main() {
	// get args, first arg is the path for the rootPackage, slice the rest to pass to go test
	packageRoot, args, err := getArgs(os.Args)
	if err != nil {
		log.Print(red(err.Error()))
		os.Exit(1)
	}

	// check results
	os.Exit(Run(packageRoot, args))
}
