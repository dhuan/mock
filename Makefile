build:
	go build -o mock
	mkdir -p bin
	mv ./mock ./bin/.

test_unit:
	go test -v $(shell git ls-files | grep _test.go)

test_e2e:
	go test -v $(shell find tests | grep _test.go)
