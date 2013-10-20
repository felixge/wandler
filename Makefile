test:
	go test -i github.com/felixge/wandler/test/integration
	go test github.com/felixge/wandler/test/integration

run: wandler-server
	@./bin/$^

wandler-server:
	@go install github.com/felixge/wandler/cmd/$@

.PHONY: test run wandler-server
