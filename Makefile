export GOPATH=${PWD}
export GOBIN=${GOPATH}/bin
export GOARCH=arm
export GOOS=linux
export GOARM=7

all:
	go get -d garage
	go build garage

clean:
	go clean garage


