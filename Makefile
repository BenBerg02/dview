BIN := build/dview
PKG := ./cmd/dview

# static cross-compiled binaries (the sqlite driver is pure Go, no CGO needed)
CROSS := CGO_ENABLED=0 GOOS=linux

.PHONY: all build build-linux-arm64 build-linux-amd64 run test vet tidy clean

all: build

build:
	go build -o $(BIN) $(PKG)

build-linux-arm64:
	$(CROSS) GOARCH=arm64 go build -o build/dview-linux-arm64 $(PKG)

build-linux-amd64:
	$(CROSS) GOARCH=amd64 go build -o build/dview-linux-amd64 $(PKG)

# usage: make run DB=path/to/file.sqlite
run: build
	./$(BIN) $(DB)

test:
	go test ./...

vet:
	go vet ./...

tidy:
	go mod tidy

clean:
	rm -rf build
