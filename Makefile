PHONY: build run run-ebiten run-ncurses clean-modules clean-all

default: build run-ncurses

clean-modules:
	@echo "Cleaning Go module cache and updating modules"
	@go clean -modcache
	@go mod tidy

clean-all: clean-modules
	@echo "Cleaning all build artifacts"
	@rm -f gol
	@go clean -cache

build:
	@echo "Building Evented GOL"
	@go build -o gol ./cmd/main.go

build-clean: clean-all
	@echo "Building Evented GOL with clean cache"
	@go build -o gol ./cmd/main.go

run-ncurses: build
	@echo "Running GOL with ncurses renderer"
	@./gol --renderer=ncurses $(ARGS)

run-ebiten: build
	@echo "Running GOL with Ebiten renderer"
	@[ -f gol ] || $(MAKE) build
	@./gol --renderer=ebiten $(ARGS)

run: run-ncurses
