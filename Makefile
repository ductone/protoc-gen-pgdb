MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: build
build:
	mkdir -p build
	go build -mod=vendor -v -o build/ . 

.PHONY: generate
generate:
	buf generate --path proto

.PHONY: example
example: build
	buf generate --template buf.example.gen.yaml --path example/models

.PHONY: fmt
fmt:
	buf format -w 

.PHONY: adddep
adddep:
	go mod tidy -v
	go mod vendor

.PHONY: updatedeps
updatedeps:
	go get -d -u ./...
	go mod tidy -v
	go mod vendor