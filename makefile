SRC_ALL := $(wildcard *.go)
SRC     := $(filter-out %_test.go,$(SRC_ALL))

run:
	go run $(SRC)

build:
	go build

test:
	go test -v

setup:
	go install github.com/datsuns/autocmd@latest
	go get -u github.com/gizak/termui/v3

auto:
	autocmd -v -t '.*\.go' -t makefile -- make test

dbg:
	@echo $(SRC_ALL)
	@echo $(SRC)

.PHONY: run test auto dbg
