BINARY='wonder.wasm'
all: build

build:
	GOOS=js GOARCH=wasm go build -o ${BINARY} ./cmd/wonder

build-prod:
	GOOS=js GOARCH=wasm go build -o ${BINARY} -ldflags "-s -w" ./cmd/wonder

.PHONY: vendor
vendor:
	GOFLAGS=-mod=vendor go mod vendor

.PHONY: serve
serve:
	go run scripts/server.go