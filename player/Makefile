.DEFAULT_GOAL := build

outfile ?= wordle-player

fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	golint ./...
.PHONY:lint

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build
.PHONY:build

clean:
	-go clean -cache -i -r
.PHONY:clean

dep: clean
	go mod download
.PHONY:dep

deploy:
	go build -o ${outfile}
.PHONY:deploy
