.PHONY: build run test

build:
	go build

run:
	python3 typo.py

test:
	cd hal && go test