##
## If you make this makefile more complex, you will learn--painfully--the
## the meaning of craniodefenestration.  
##
## This should just call the go commands, nothing else.  It should be 
## run from the src/github.com/iansmith/movienight directory using the 
## installed Godeps.  Really, don't add anything complicated; don't.
##
SHA:=$(shell git rev-parse HEAD)

all: build

.PHONY: build
build:
	godep go install  -ldflags '-X $(GO_GIT_DESCRIBE_SYMBOL) $(SHA)' github.com/iansmith/movienight/server
	godep go install github.com/iansmith/movienight/migrate
	godep go install github.com/iansmith/movienight/tooling/pagegen

.PHONY: godeps
godeps:
	rm -rf Godeps
	godep save ./...
