# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	@go build -o ./tmp/main ./main.go

# Run the application
run:
	@go run main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching...";\
	else \
		read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching...";\
		else \
			echo "You chose not to install air. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Lint the application
lint:
	@echo "Running linter..."
	@golangci-lint run --fix --verbose

# Align the structs
align:
	@echo "Aligning structs..."
	@if command -v betteralign > /dev/null; then \
		betteralign -apply ./...; \
	else \
		read -p "Go's 'betteralign' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/dkorunic/betteralign@latest; \
			betteralign -apply ./...; \
		else \
			echo "You chose not to install betteralign. Exiting..."; \
			exit 1; \
		fi; \
	fi

# Generate swagger documentation
swagger:
	@echo "Generating swagger..."
	@if command -v swag > /dev/null; then \
		swag init --parseDependency --parseInternal --generatedTime -g ./main.go -o ./docs/swagger; \
	else \
		read -p "Go's 'swag' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/swaggo/swag/cmd/swag@latest; \
		swag init --parseDependency --parseInternal --generatedTime -g ./main.go -o ./docs/swagger; \
		else \
			echo "You chose not to install swag. Exiting..."; \
			exit 1; \
		fi; \
	fi

proto:
	@echo "Generating proto..."
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./internal/proto/*.proto

templ:
	@echo "Generating templates..."
	@if command -v templ > /dev/null; then \
		templ generate; \
		tailwindcss -i ./static/css/input.css -o ./static/css/output.css; \
	else \
		read -p "Go's 'templ' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			templ generate; \
		else \
			echo "You chose not to install templ. Exiting..."; \
			exit 1; \
		fi; \
	fi

tailwind:
	@echo "Generating tailwind..."
	@if command -v tailwindcss > /dev/null; then \
		tailwindcss -i ./static/css/input.css -o ./static/css/output.css; \
	else \
		read -p "Go's 'tailwind' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64; \
			chmod +x tailwindcss-macos-arm64; \
			mv tailwindcss-macos-arm64 tailwindcss; \
			tailwindcss -i ./static/css/input.css -o ./static/css/output.css; \
		else \
			echo "You chose not to install tailwind. Exiting..."; \
			exit 1; \
		fi; \
	fi

.PHONY: all build run test clean watch itest lint align swagger proto templ tailwind
