.PHONY: all build-local build-cloud run-local run-cloud test clean

all: build-local build-cloud

build-local:
	@echo "Building local backend..."
	CGO_ENABLED=1 go build -o bin/local-backend ./cmd/local/main.go

build-cloud:
	@echo "Building cloud backend..."
	CGO_ENABLED=0 GOOS=linux go build -o bin/cloud-backend ./cmd/cloud/main.go

run-local:
	@echo "Running local backend..."
	go run ./cmd/local/main.go

run-cloud:
	@echo "Running cloud backend..."
	go run ./cmd/cloud/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f shosha_finance.db

deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

lint:
	@echo "Running linter..."
	golangci-lint run

build-windows:
	@echo "Building for Windows..."
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o bin/local-backend.exe ./cmd/local/main.go

build-darwin:
	@echo "Building for macOS..."
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/local-backend-darwin ./cmd/local/main.go
