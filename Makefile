BINARY     = scope
MODULE     = github.com/lucap123/scope
VERSION   ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS    = -ldflags "-X $(MODULE)/cmd.Version=$(VERSION) -s -w"
BUILD_DIR  = dist

.PHONY: build install test clean release

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

test:
	go test ./... -v -count=1

clean:
	rm -rf $(BINARY) $(BUILD_DIR)

release: clean
	mkdir -p $(BUILD_DIR)
	GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-amd64   .
	GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-linux-arm64   .
	GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-amd64  .
	GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-darwin-arm64  .
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY)-windows-amd64.exe .
	@echo "Binaries built in $(BUILD_DIR)/"

fmt:
	gofmt -w .

vet:
	go vet ./...
