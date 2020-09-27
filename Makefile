.PHONY: all
all:

.PHONY: test
test:
	@go test -count 1 ./...
