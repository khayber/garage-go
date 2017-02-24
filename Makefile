export GOPATH=${PWD}

all: get build install

pi0: export GOARCH=arm
pi0: export GOOS=linux
pi0: export GOARM=5
pi0: build install

pi3: export GOARCH=arm
pi3: export GOOS=linux
pi3: export GOARM=7
pi3: build install

get:
	go get garage

build:
	go build garage

install:
	go install garage

clean:
	go clean garage

purge: clean
	rm garage
	rm -rf bin/*
	rm -rf pkg/*
