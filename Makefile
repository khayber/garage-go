export GOPATH=${PWD}

all: get build

pi0: export GOARCH=arm
pi0: export GOOS=linux
pi0: export GOARM=5
pi0: get build

pi3: export GOARCH=arm
pi3: export GOOS=linux
pi3: export GOARM=7
pi3: get build

get:
	go get -d

build:
	go build

clean:
	go clean

purge: clean
	-rm -rf src 
