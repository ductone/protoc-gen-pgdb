MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables

.PHONY: install
install:
	go install -mod=vendor -v .

.PHONY: generate
generate:
	buf generate

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