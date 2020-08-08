
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BUILD_PATH=./bin
BINARY_NAME=aws-chaos-cli
BINARY_UNIX=$(BINARY_NAME)-linux-amd64
BINARY_WINDOWS=$(BINARY_NAME)-windows-amd64
BINARY_MACOS=$(BINARY_NAME)-darwin-amd64

all: test build
build:
				$(GOBUILD) -o $(BINARY_NAME) -v
test:
				$(GOTEST) -v ./...
clean:
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
				rm -f $(BINARY_UNIX)
				rm -f $(BINARY_MACOS)
run:
				$(GOCLEAN)
				$(GOBUILD) -o $(BINARY_NAME) -v
				./$(BINARY_NAME)
deps:
				$(GOGET) github.com/spf13/cobra/cobra
				$(GOGET) ./...


# Cross compilation
build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_PATH)/$(BINARY_UNIX) -v
build-windows:
				CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_PATH)/$(BINARY_WINDOWS) -v
build-macos:
				CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_PATH)/$(BINARY_MACOS) -v
docker-build:
				docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/patmizi/aws-chaos-cli golang:latest go build -o "$(BINARY_UNIX)" -v