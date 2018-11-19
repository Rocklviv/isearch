# GO parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Application related variables
APP_BINARY_NAME=isearch
APP_BINARY_DIR=release

# Multiarch build related variables
UNIX_BINARY_DIR=linux
OSX_BINARY_DIR=darwin
WINDOWS_BINARY_DIR=win64

all: lint test build

build:
	$(GOBUILD) -o $(APP_BINARY_DIR)/$(APP_BINARY_NAME) -v

test:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -func=coverage.out
	$(GOCMD) tool cover -html=coverage.out

lint:
	$(GOCMD) fmt ./...

clean:
	$(GOCLEAN)
	rm -rf $(APP_BINARY_DIR)

deps:
	$(GOGET) github.com/olekukonko/tablewriter
	$(GOGET) github.com/stretchr/testify

build-multiarch:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(APP_BINARY_DIR)/$(OSX_BINARY_DIR)/$(APP_BINARY_NAME) -v
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(APP_BINARY_DIR)/$(UNIX_BINARY_DIR)/$(APP_BINARY_NAME) -v
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(APP_BINARY_DIR)/$(WINDOWS_BINARY_DIR)/$(APP_BINARY_NAME) -v