.PHONY: build
build:
	go build cmd/omp-demo-api/main.go

.PHONY: test
test:
	go test -v ./...