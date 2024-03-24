SRC_ALL := $(wildcard *.go)
SRC     := $(filter-out %_test.go,$(SRC_ALL))

run:
	go run -race $(SRC)

build:
	go build 

release:
	go build \
		-a \
		-ldflags="-s -w"


test:
	go test -v

setup:
	go install github.com/datsuns/autocmd@latest
	go install github.com/go-delve/delve/cmd/dlv@latest
	go get -u github.com/rivo/tview
	go get -u github.com/pelletier/go-toml/v2
	go get -u github.com/cockroachdb/errors

auto:
	autocmd -v -t '.*\.go' -t makefile -- make test

gdb:
	dlv debug

dbg:
	@echo $(SRC_ALL)
	@echo $(SRC)

.PHONY: run build test setup auto dbg
