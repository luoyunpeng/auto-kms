# Go parameters
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get
BINARY_NAME=kms
BINARY_UNIX=$(BINARY_NAME)_unix
Version=1.0.0

all: test build
build: vendor
	$(GOBUILD) -ldflags="-w -s" -v -o $(BINARY_NAME) ./cmd/kms/kms.go
build-v:
	$(GOBUILD) -ldflags="-w -s" -v -mod=vendor -o $(BINARY_NAME) ./cmd/kms/kms.go
vendor:
	go env -w GOPROXY=https://goproxy.cn,direct
rm-build:
	rm -f $(BINARY_NAME)
test:
	# TODO, add test, $(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run: build-v
	cp -f ./asset/config/block.yml .
	./$(BINARY_NAME)
	#./$(BINARY_NAME)
nohup: build-v
	cp -f ./asset/config/block.yml .
	nohup ./$(BINARY_NAME) & echo "$$!" > block.pid
stop:
	kill $(shell cat block.pid)
rerun: stop nohup

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(PWD)":/opt/git/blockMGR/ -w /opt/git/blockMGR/ golang:1.15.2-alpine go build -v -mod=vendor -o "$(BINARY_NAME)" ./cmd/kms/kms.go
image-build:
