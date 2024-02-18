
buildcli:
	@go build -o bin/cli cmd/cli/main.go

buildserver:
	@go build -o bin/server cmd/server/main.go

devapi:
	@gow -c run apps/api/main.go
