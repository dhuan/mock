build:
	go build -o mock
	mkdir -p bin
	mv ./mock ./bin/.

test: test_unit test_e2e

test_unit:
	go test -v ./internal/... ./pkg/...

test_e2e:
	go test -v $(shell find tests | grep 'e2e.*_test.go')

doc_dev:
	sh ./scripts/doc_dev.sh

doc_build:
	sh ./scripts/doc_build.sh

github_doc_vars:
	sh ./scripts/github_doc_vars.sh
