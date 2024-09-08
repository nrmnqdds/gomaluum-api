dev:
	go run ./cmd/main.go

benchmark:
	cd test && go test -bench=.

lint:
	golangci-lint run

swagger:
	swag init --parseDependency --parseInternal -g ./cmd/main.go -o ./docs/swagger
