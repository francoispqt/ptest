
[![Build Status](https://travis-ci.org/francoispqt/ptest.svg?branch=master)](https://travis-ci.org/francoispqt/ptest)
[![codecov](https://codecov.io/gh/francoispqt/ptest/branch/master/graph/badge.svg)](https://codecov.io/gh/francoispqt/ptest)
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
