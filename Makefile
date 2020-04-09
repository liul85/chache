default: test

test:
		go test ./... -v

.PHONY: test