 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_DIR=/opt
BUILDPARAMS="-a -installsuffix nocgo"

all: test build
build:
	 $(GOBUILD) $(BUILDPARAMS) -o $(BINARY_NAME) -v ./...
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	dep ensure -update


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux $(GOBUILD) $(BUILDPARAMS) -o $(BINARY_DIR) -v ./...
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/fedjo/ori-app golang:1.13-alpine go build -o "$(BINARY_DIR)" -v ./...

