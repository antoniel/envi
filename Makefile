.PHONY: start dev

start:
	go build -o myapp
	./myapp

dev:
	@go run cmd/main.go
