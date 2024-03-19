SRC_ALL := $(wildcard *.go)
SRC     := $(filter-out %_test.go,$(SRC_ALL))

run:
	go run $(SRC)

test:
	go test -v

setup:

auto:
	autocmd -v -t '.*\.go' -t makefile -- make test

dbg:
	@echo $(SRC_ALL)
	@echo $(SRC)

.PHONY: run test auto dbg
