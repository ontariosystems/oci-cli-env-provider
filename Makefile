.PHONY: lint build prep test watch cover
.DEFAULT_GOAL := test

lint:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest run

test:
	go run github.com/onsi/ginkgo/v2/ginkgo -r -cover

watch:
	go run github.com/onsi/ginkgo/v2/ginkgo watch -r -cover

cover:
	go tool cover -html=coverprofile.out
