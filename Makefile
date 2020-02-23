.PHONY: mkdir-build
mkdir-build:
	if [ ! -d build ]; then mkdir build; fi

.PHONY: test-cli
test-cli:
	cd cmd/httpmon && go test ./...

.PHONY: test
test: mkdir-build
	go test -coverprofile=build/test-report.txt
	go tool cover -html=build/test-report.txt -o build/test-report.html

.PHONY: clean
clean:
	go clean
	rm -rf build
	cd cmd/httpmon && go clean

.PHONY: check
check: test test-cli
