
.PHONY: build
build:
	go build -o bin/shopcart cmd/main.go

.PHONY: build
tidy:
	go mod tidy
