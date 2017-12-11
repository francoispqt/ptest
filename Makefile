.PHONY: test
test: 
	go test

.PHONY: cover
cover: 
	go test -coverprofile=coverage.out

.PHONY: coverhtml
coverhtml: 
	 go tool cover -html=coverage.out