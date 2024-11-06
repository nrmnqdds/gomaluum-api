dev:
	./dev.sh

benchmark:
	cd test && go test -bench=.

lint:
	golangci-lint run --fix --verbose

swagger:
	swag init --parseDependency --parseInternal -g ./cmd/main.go -o ./docs/swagger

update:
	go get -u ./...
	go mod tidy

align:
	# go get -u https://github.com/dkorunic/betteralign
	betteralign -apply ./...
