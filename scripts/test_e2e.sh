export PATH=$PATH:$(pwd)/bin

go test -v $(find tests | grep 'e2e.*_test.go')
