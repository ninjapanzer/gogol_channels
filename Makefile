PHONY: build run

default: build run

build:
	@echo "Building Evented GOL"
	@go vet ./cmd/main.go
	@go build -o gol ./cmd/main.go

run:
	@echo "running GOL"
	@./gol
