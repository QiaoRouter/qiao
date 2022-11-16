.PHONY: build run test

build:
	go build

run:
	go build && ./qiao

test-hal:
	cd hal && go test