.PHONY: build-cli-mac
build-cli-mac:
	cd cmd/httpmon && GOOS=darwin GOARCH=amd64 go build -o ../../build/mac/http-mon ./...

.PHONY: build-cli-win
build-cli-win:
	cd cmd/httpmon && GOOS=windows GOARCH=amd64 go build -o ../../build/win/http-mon.exe ./...

.PHONY: build-cli-linux
build-cli-linux:
	cd cmd/httpmon && GOOS=linux GOARCH=amd64 go build -o ../../build/linux/http-mon ./...

.PHONY: build
build: build-cli-linux build-cli-mac build-cli-win
	mkdir build/release
	zip build/release/http-mon-darwin-amd64.zip build/mac/http-mon
	zip build/release/http-mon-win-amd64.zip build/win/http-mon.exe
	zip build/release/http-mon-linux-amd64.zip build/linux/http-mon

.PHONY: test-cli
test-cli:
	cd cmd/httpmon && go test ./...

.PHONY: test
test:
	go test ./...

.PHONY: tidy-root
tidy-root:
	go mod tidy

.PHONY: tidy-cli
tidy-cli:
	cd cmd/httpmon && go mod tidy

.PHONY: tidy
tidy: tidy-root tidy-cli

.PHONY: clean
clean:
	go clean
	rm -rf build
	cd cmd/httpmon && go clean

.PHONY: check
check: test test-cli
