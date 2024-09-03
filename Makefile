# Go parameters
GOCMD=go
GOBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

test:
	go run -mod=vendor ./cmd/app/balance_monitor.go -config ./depolyments/config-sample.toml

install:
	go mod mode

build:
	go build -mod=vendor -v -o ./bin/balance_monitor ./cmd/app 

build_linux:
	$(GOBUILD) -mod=vendor -v -o ./bin/balance_monitor ./cmd/app