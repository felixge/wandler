test:
	@./test.bash

run: wandler-server
	@./bin/$^

wandler-server:
	@go install github.com/felixge/wandler/cmd/$@

.PHONY: test run wandler-server
