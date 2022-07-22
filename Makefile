build:
	go build -o mock
	mkdir -p bin
	mv ./mock ./bin/.

test: test_unit test_e2e

test_unit:
	go test -v $(shell git ls-files | grep _test.go | grep -v e2e)

test_e2e:
	go test -v $(shell find tests | grep 'e2e.*_test.go')
