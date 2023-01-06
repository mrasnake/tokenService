build:
	go build -o ./cmd/api ./cmd/api

run: build
	./cmd/api/api

help: build
	./cmd/api/api -h
