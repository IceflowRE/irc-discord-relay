# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=irc-discord-relay

all: deps clean build
build:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -x
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -rf ./build/
run:
	$(GOBUILD) -o ./build/$(BINARY_NAME) -x
	./build/$(BINARY_NAME)
deps:
	$(GOGET) ./...
