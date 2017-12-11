// Package ptest is small package used to run tests in package and subpackages concurrently
// adds a slightly better output
package main

import (
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
	vendorDir = "vendor"
)

// variables, goPath, colors and al
var (
	goPath = os.Getenv("GOPATH")
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

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
	return
}

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
func NewTests(dirPath string, testsResultChan chan TestResult, args []string) []TestRunner {
	tests := make([]TestRunner, 0)
	files, err := ioutil.ReadDir(goPath + "/src/" + dirPath)
	if err != nil {
		panic(err)
	}
	hasTest := false
	for _, f := range files {
		// avoid testing
		if f.IsDir() && f.Name() != vendorDir {
			// should run concurrently
			tests = append(tests, NewTests(dirPath+"/"+f.Name(), testsResultChan, args)...)
		} else if strings.HasSuffix(f.Name(), "_test.go") && hasTest == false {
			test := Test{dirPath}
			tests = append(tests, test)
			// run test
			go test.Run(testsResultChan, args)
			hasTest = true
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
			log.Printf(green("✓ %s"), string(testResult.Package))
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

func getArgs(args []string) (string, []string) {
	return args[1], args[2:]
}

func main() {
	// get args, first arg is the path for the rootPackage, slice the rest to pass to go test
	packageRoot, args := getArgs(os.Args)
	// create tests by recursively getting subpackages
	testsResultChan := make(chan TestResult)
	// get tests and run
	tests := NewTests(packageRoot, testsResultChan, args)
	// check results
	os.Exit(GetTestsResult(testsResultChan, tests))
}
