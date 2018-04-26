
[![Build Status](https://travis-ci.org/francoispqt/ptest.svg?branch=master)](https://travis-ci.org/francoispqt/ptest)
[![codecov](https://codecov.io/gh/francoispqt/ptest/branch/master/graph/badge.svg)](https://codecov.io/gh/francoispqt/ptest)
[![Go Report Card](https://goreportcard.com/badge/github.com/francoispqt/ptest)](https://goreportcard.com/report/github.com/francoispqt/ptest)

# Ptest
Proxy to run go tests for a package and all its subpackage in parallel for faster testing

## Usage
Get the package
```shell
$ go install github.com/francoispqt/ptest
```
Run it for a specific package
```shell
$ ptest github.com/francoispqt/ptest
```
Remaining parameters are passed to go test
```shell
$ ptest github.com/francoispqt/ptest -race -coverprofile=profile.out -covermode=atomic
```

## Skip test for a package (but not its subpackages)
You can skip a package by creating a file named `.ptestskip` in the root directory of the package

## Output
Ptest has a slightly better output than the regular `go test` (it has colours :D)

![alt text][screenshot]

[screenshot]: https://preview.ibb.co/iefuHw/Screen_Shot_2017_12_12_at_9_29_39_am.png "Screenshot ptest"