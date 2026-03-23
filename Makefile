.PHONY: build clean install run

BINARY_NAME=hati
VERSION=0.1.0
GO_LDFLAGS=-ldflags "-X main.version=${VERSION}"

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

BUILD_DIR=./bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(GO_LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/hati

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

install: build
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
