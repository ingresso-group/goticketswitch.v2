GOPRIVATE = github.com/ingresso-group
CGO_ENABLED=0
GOOS=linux

.PHONY: setup
setup:
	go install github.com/kyoh86/richgo@latest
	go install github.com/jstemmer/go-junit-report@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: fmt
fmt:
	goimports -w .
	gofmt -s -w .

.PHONY: govendor
govendor:
	go mod tidy -compat=1.19
	go mod vendor

.PHONY: goupgrade
goupgrade:
	GOPRIVATE=$(GOPRIVATE) go get -u -v ./...
	$(MAKE) govendor

.PHONY: test
test:
	richgo test -v -race -coverpkg=./... -coverprofile=coverage.txt ./... -mod=readonly

.PHONY: golangci-lint
golangci-lint:
	golangci-lint run

.PHONY: lint
lint: fmt golangci-lint
