.PHONY: build run test

build:
	go build

run:
	go build && python3 typo.py

test:
	cd hal && go test

run-r1:
	bird6 -c ./setup/ripng/netns/v1/bird-r1.conf -d -s ~/bird-r1.ctl

run-r3:
	bird6 -c ./setup/ripng/netns/v1/bird-r3.conf -d -s ~/bird-r3.ctl

perf:
	go build && python3 perf_qiao.py