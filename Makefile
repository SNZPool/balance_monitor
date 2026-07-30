# Go parameters
GOCMD=go
GOBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

.PHONY: test install build build_linux clean

test:
	go run ./cmd/app/balance_monitor.go -config ./deployments/config-sample.json

install:
	@GOPROXY=https://proxy.golang.org,direct go mod tidy

build:
	mkdir -p ./bin
	go build -v -o ./bin/balance_monitor ./cmd/app

# Static linux/amd64 binary for production hosts
build_linux:
	mkdir -p ./bin
	$(GOBUILD) -ldflags="-s -w" -v -o ./bin/balance_monitor ./cmd/app

clean:
	rm -rf ./bin
