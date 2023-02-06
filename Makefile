.PHONY: clean

all: eve eve-tools

# Build executable for Eve program
eve:
	go mod download
	go build --ldflags "-s -w" -o bin/eve ./cmd/eve/

eve-tools:
	go mod download
	go build --ldflags "-s -w" -o bin/eve-tools ./cmd/eve-tools/

unit-test:
	go test ./... -v -race -coverprofile=coverage.out -covermode=atomic -count=1

integration-test:
	./bin/eve --log-format pretty &
	sleep 5
	go test ./... -v --tags=integration -count=1
	killall eve

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
