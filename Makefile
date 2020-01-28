.PHONY: build-mon-mac
build-mon-mac:
	cd cmd/httpmon && go build -o ../../build/http-mon ./...

.PHONY: test
test:
	go test ./...
