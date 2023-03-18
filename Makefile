build:
	go build -o mock
	mkdir -p bin
	mv ./mock ./bin/.

test: test_unit test_e2e

test_unit: test_unit_internal test_unit_pkg

test_unit_internal:
	find . | grep _test.go | grep internal | xargs -I _FILE_ go test -v _FILE_

test_unit_pkg:
	go test -v $(shell git ls-files | grep 'pkg/' | grep _test.go | grep -v e2e)

test_e2e:
	go test -v $(shell find tests | grep 'e2e.*_test.go') -run "Test_E2E_Response_Exec"
