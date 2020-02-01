.PHONY: build-cli-mac
build-cli-mac:
	cd cmd/httpmon && GOOS=darwin GOARCH=amd64 go build -o ../../build/mac/http-mon ./...

.PHONY: build-cli-win
build-cli-win:
	cd cmd/httpmon && GOOS=windows GOARCH=amd64 go build -o ../../build/win/http-mon ./...

.PHONY: build-cli-linux
build-cli-linux:
	cd cmd/httpmon && GOOS=linux GOARCH=amd64 go build -o ../../build/linux/http-mon ./...

.PHONY: build
build: build-cli-linux build-cli-mac build-cli-win

.PHONY: test-cli
test-cli:
	cd cmd/httpmon && go test ./...

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	go clean
	rm -rf build
	cd cmd/httpmon && go clean

.PHONY: check
check: test test-cli
