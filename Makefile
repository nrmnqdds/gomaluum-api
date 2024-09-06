dev:
	go run ./cmd/main.go

benchmark:
	cd test && go test -bench=.

lint:
	golangci-lint run
