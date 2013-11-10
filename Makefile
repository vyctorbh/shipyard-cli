all:
	@go get github.com/ehazlett/shipyard-go/shipyard
	@go get github.com/codegangsta/cli
	@cd ./cli && go build -o ../shipyard

fmt:
	@cd ./cli && go fmt
clean:
	@rm -rf shipyard
