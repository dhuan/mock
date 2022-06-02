build:
	go build -o mock
	mkdir -p bin
	mv ./mock ./bin/.

test_unit:
	go test -v $(shell git ls-files | grep _test.go)

test: test_unit
