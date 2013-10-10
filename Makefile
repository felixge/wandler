run: wandler-server
	@./bin/$^

wandler-server:
	@go install github.com/felixge/wandler/cmd/$@

.PHONY: run wandler-server
