.PHONY: build
#: Performs a clean run of the project
build: deps
	@go build -o strangelet.exe cmd/main.go

.PHONY: run
#: Starts the project
run: build
	@./strangelet.exe

.PHONY: clean
#: Cleans slate for docker
clean:
	@rm strangelet.exe

.PHONY: deps
#: Install dependencies for docker and targets in this makefile
deps:
	@go mod tidy -compat=1.17

.PHONY: help
#: Lists available commands
help:
	@echo "Available Commands for project:"
	@grep -B1 -E "^[a-zA-Z0-9_-]+\:([^\=]|$$)" Makefile \
	 | grep -v -- -- \
	 | sed 'N;s/\n/###/' \
	 | sed -n 's/^#: \(.*\)###\(.*\):.*/\2###\1/p' \
	 | column -t  -s '###'
