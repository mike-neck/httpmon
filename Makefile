.PHONY: mkdir-build
mkdir-build:
	if [ ! -d build ]; then mkdir build; fi

.PHONY: test-gen
test-gen:
	${GOPATH}/bin/mockgen -source=client.go -package=httpmon -destination=./client_mock.go

.PHONY: test-cli
test-cli:
	cd cmd/httpmon && go test ./...

.PHONY: test
test: mkdir-build test-gen
	go test -coverprofile=build/test-report.txt
	go tool cover -html=build/test-report.txt -o build/test-report.html

.PHONY: coverage
coverage:
	go tool cover -html=build/test-report.txt -o build/test-report.html

.PHONY: clean
clean:
	go clean
	rm -rf build
	if [ -f ./client_mock.go ]; then rm ./client_mock.go ; fi
	cd cmd/httpmon && go clean

.PHONY: check
check: test test-cli
