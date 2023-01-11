.PHONY: clean

all: eve eve-tools

# Build executable for Eve program
eve:
	go mod download
	go build --ldflags "-s -w" -o bin/eve ./cmd/eve/

eve-tools:
	go mod download
	go build --ldflags "-s -w" -o bin/eve-tools ./cmd/eve-tools/

test:
	go clean -testcache
	go test ./... -v -race -coverprofile=coverage.out -covermode=atomic

# Build and execute Eve program
start: eve
	./bin/eve --log-format pretty

# Format Sojourner source code with Go toolchain
format:
	go mod tidy
	go fmt ./...

# Clean up binary output folder
clean:
	rm -rf bin/
