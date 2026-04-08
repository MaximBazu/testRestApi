.PHONY: build
build:
	go build -o ./bin/myapp/myapp.exe cmd/myapp/main.go

.PHONY: run
run: build
	./bin/myapp/myapp.exe