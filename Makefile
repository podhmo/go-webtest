.PHONY: lint

test:
	go test ./...

lint:
	golangci-lint run # -v --enable-all

format:
	goimports -w $(shell git ls-files | grep ".go$$")

# example (for testing)
httpbin:
	cd examples/00httpbin && go run ./cmd/httpbin/main.go localhost:8080
.PHONY: httpbin
