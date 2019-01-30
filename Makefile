.PHONY: lint

lint:
	golangci-lint run -v #--enable-all
