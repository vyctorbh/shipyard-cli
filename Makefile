all:
	@go get -d -v ./...
	@cd ./cli && go build -o ../shipyard

fmt:
	@go fmt ./...

test:
	@go test ./...
clean:
	@rm -rf shipyard
