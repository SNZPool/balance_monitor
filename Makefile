# Go parameters
GOCMD=go
GOBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

test:
	go run ./cmd/app/balance_monitor.go -config ./depolyments/config-sample.toml

install:
	@GOPROXY=https://mycompany.com/proxy,direct go mod tidy

build:
	go build -v -o ./bin/balance_monitor ./cmd/app 

build_linux:
	$(GOBUILD) -v -o ./bin/balance_monitor ./cmd/app