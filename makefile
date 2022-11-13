.PHONY: build run test

build:
	go build

run:
	go build && ./qiao

test-hal:
	cd hal && go test

test-hal-R2:
	cd hal && ip netns exec R2 go test