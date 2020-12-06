# Go parameters
GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get
BINARY_NAME=block
BINARY_UNIX=$(BINARY_NAME)_unix
Version=1.0.0

all: test build
build: vendor
	$(GOBUILD) -ldflags="-w -s" -v -o $(BINARY_NAME) ./cmd/block-data/block_data.go
build-tmp: build rm-build
build-dev: vendor
	$(GOBUILD) -ldflags="-w -s" -v -o $(BINARY_NAME)_dev ./cmd/block-data/block_data.go
build-test: vendor
	$(GOBUILD) -ldflags="-w -s" -v -o $(BINARY_NAME)_test ./cmd/block-data/block_data.go
build-v:
	$(GOBUILD) -ldflags="-w -s" -v -mod=vendor -o $(BINARY_NAME) ./cmd/block-data/block_data.go
build-dev-v:
	$(GOBUILD) -ldflags="-w -s" -v -mod=vendor -o $(BINARY_NAME)_dev ./cmd/block-data/block_data.go
build-test-v:
	$(GOBUILD) -ldflags="-w -s" -v -mod=vendor -o $(BINARY_NAME)_test ./cmd/block-data/block_data.go
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
nohup-dev: build-dev-v
	cp -f ./asset/config/block-dev.yml ./block.yml
	nohup ./$(BINARY_NAME)_dev & echo "$$!" > block-dev.pid
nohup-test: build-test-v
	cp -f ./asset/config/block-test.yml ./block.yml
	nohup ./$(BINARY_NAME)_test & echo "$$!" > block-test.pid
stop:
	kill $(shell cat block.pid)
stop-dev:
	kill $(shell cat block-dev.pid)
stop-test:
	kill $(shell cat block-test.pid)
rerun: stop nohup
rerun-dev: stop-dev nohup-dev
rerun-test: stop-test nohup-test

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	#docker run --rm -it -v "$(GOPATH)":/go -w /go/src/github.com/luoyunpeng/monitor golang:1.12.6-alpine go build -v -tags=jsoniter -o "$(BINARY_NAME)" ./cmd/monitor/monitor.go
	docker run --rm -it -v "$(PWD)":/opt/git/blockMGR/ -w /opt/git/blockMGR/ golang:1.15.2-alpine go build -v -mod=vendor -o "$(BINARY_NAME)" ./cmd/block-data/block_data.go
image-build: docker-build
	mkdir block-mgr && mv $(BINARY_NAME) block-mgr/
	cp asset/config/block-template.yml block-mgr/block.yml
	cp asset/shell/entrypoint.sh block-mgr/
	docker build -t block-mgr:$(Version) .
	rm -rf block-mgr
# docker container run -d --name block-mgr
# -e MQTT_SERVER_URI=139.159.225.207:8083 -e MQTT_USER=status-mgr -e MQTT_PASSWORD=test
# -e MYSQL_SERVER_IP=222.85.139.245 -e MYSQL_PORT=18306 -e MYSQL_USER=root -e MYSQL_PASSWORD=123
# -e MONGODB_SERVER_IP=192.168.100.210 -e MONGODB_PORT=1666 -e MONGODB_USER=root -e MONGODB_PASSWORD=iscas123
# -p 12345:12345--restart=always block-mgr:$ver
