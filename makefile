
buildcli:
	@go build -o bin/cli cmd/cli/main.go

buildserver:
	@go build -o bin/server cmd/server/main.go

api:
	@gow -c run ./apps/api 

cli:
	@go run ./apps/envi-cli/main.go
