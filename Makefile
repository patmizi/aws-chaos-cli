
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=aws-chaos-cli
BINARY_UNIX=$(BINARY_NAME)-linux
BINARY_WINDOWS=$(BINARY_NAME)-windows

all: test build
build:
				$(GOBUILD) -o $(BINARY_NAME) -v
test:
				$(GOTEST) -v ./...
clean:
				$(GOCLEAN)
				rm -f $(BINARY_NAME)
				rm -f $(BINARY_UNIX)
run:
				$(GOCLEAN)
				$(GOBUILD) -o $(BINARY_NAME) -v
				./$(BINARY_NAME)
deps:
				$(GOGET) github.com/spf13/cobra/cobra


# Cross compilation
build-linux:
				CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX)-amd64 -v
build-windows:
				CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS)-amd64 -v
docker-build:
				docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/patmizi/aws-chaos-cli golang:latest go build -o "$(BINARY_UNIX)" -v