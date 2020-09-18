build:
	@cd cmd/server && go build -o ../../bin/server

test:
	go test -race ./... -v

fmt:
	go fmt ./...
