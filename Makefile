.PHONY: mkdir-build
mkdir-build:
	if [ ! -d build ]; then mkdir build; fi

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
