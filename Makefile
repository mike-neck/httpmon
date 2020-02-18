.PHONY: build-cli-mac
build-cli-mac:
	cd cmd/httpmon && GOOS=darwin GOARCH=amd64 go build -o ../../build/mac/http-mon ./...

.PHONY: build-cli-win
build-cli-win:
	cd cmd/httpmon && GOOS=windows GOARCH=amd64 go build -o ../../build/win/http-mon.exe ./...

.PHONY: build-cli-linux
build-cli-linux:
	cd cmd/httpmon && GOOS=linux GOARCH=amd64 go build -o ../../build/linux/http-mon ./...

.PHONY: mkdir-build
mkdir-build:
	if [ ! -d build ]; then mkdir build; fi

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
test: mkdir-build
	go test -coverprofile=build/v1-report.txt ./...
	go tool cover -html=build/v1-report.txt -o build/v1-report.html

.PHONY: clean
clean:
	go clean
	rm -rf build
	cd cmd/httpmon && go clean

.PHONY: check
check: test test-cli
