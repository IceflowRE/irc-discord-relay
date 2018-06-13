# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=irc-discord-relay

ifeq ($(OS),Windows_NT)
	PREFIX=.exe
else
	PREFIX=""
endif

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME)$(PREFIX) -x
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)$(PREFIX)
run:
	$(GOBUILD) -o $(BINARY_NAME)$(PREFIX) -x
	./$(BINARY_NAME)$(PREFIX)
deps:
	$(GOGET) ./...
