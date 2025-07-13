golint:
	golangci-lint run -c .golangci.yaml

golintfix:
	golangci-lint run -c .golangci.yaml --fix

gofmt:
	goimports -w server