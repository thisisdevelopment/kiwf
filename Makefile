GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: build

workdir:
	mkdir -p workdir

build: workdir/kiwf

build-native: $(GOFILES)
	go build -o workdir/native-kiwf .

workdir/kiwf: $(GOFILES)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o workdir/kiwf .

test: test-all

test-all:
	@go test -v $(GOPACKAGES)