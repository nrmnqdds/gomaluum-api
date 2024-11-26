benchmark:
	cd test && go test -bench=.

lint:
	@echo "Running linter..."
	@golangci-lint run --fix --verbose

swagger:
	@echo "Generating swagger..."
	@if command -v betteralign > /dev/null; then \
		swag init --parseDependency --parseInternal -g ./cmd/main.go -o ./docs/swagger; \
	else \
		read -p "Go's 'swag' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/swaggo/swag/cmd/swag@latest; \
			betteralign -apply ./...; \
		else \
			echo "You chose not to install betteralign. Exiting..."; \
			exit 1; \
		fi; \
	fi

update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

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

dev:
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

