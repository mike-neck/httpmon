.PHONY: mkdir-build
mkdir-build:
	if [ ! -d build ]; then mkdir build; fi

.PHONY: test-gen
test-gen:
	${GOPATH}/bin/mockgen -source=client.go -package=httpmon -destination=./client_mock.go
	${GOPATH}/bin/mockgen -source=use_case.go -package=httpmon -destination=./use_case_mock.go
	${GOPATH}/bin/mockgen -source=httpmon.go -package=httpmon -destination=./httpmon_mock.go

.PHONY: test-cli
test-cli:
	cd cmd/httpmon && go test ./...

.PHONY: test
test: mkdir-build test-gen
	go test -coverprofile=build/test-report.txt

.PHONY: coverage
coverage:
	go tool cover -html=build/test-report.txt -o build/test-report.html

.PHONY: clean
clean:
	go clean
	rm -rf build
	ls | grep _mock | xargs rm
	cd cmd/httpmon && go clean

.PHONY: check
check: test test-cli
