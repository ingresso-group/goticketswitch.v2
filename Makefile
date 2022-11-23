.PHONY: test setup

test:
	go test ./...
	go test -cover ./...
	golint ./...
	go tool vet -all .
	gocyclo -over 10 .

setup:
	go get -v ./...
	go get github.com/golang/lint/golint
	go get github.com/fzipp/gocyclo
