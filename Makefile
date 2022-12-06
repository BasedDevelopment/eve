.PHONY: clean

all: eve

# Build executable for Eve program
eve:
	go mod download
	go build --ldflags "-s -w" -o bin/eve ./cmd/eve/main.go

# Build and execute Eve program
start: eve
	./bin/eve

# Format Sojourner source code with Go toolchain
format:
	go mod tidy
	go fmt ./...

# Clean up binary output folder
clean:
	rm -rf bin/
