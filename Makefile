.PHONY: build-mon-mac
build-mon-mac:
	cd cmd/httpmon && go build -o ../../build/http-mon ./...

.PHONY: test-mon
test-mon:
	cd cmd/httpmon && go test ./...

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean
	rm -rf build
