VERSION ?= $(shell git describe --tags --always --dirty)

build:
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o treekei ./cmd/treekei/
