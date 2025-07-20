gobuild:
	go build -o server.out

golint:
	golangci-lint run -c .golangci.yaml

golintfix:
	golangci-lint run -c .golangci.yaml --fix

gofmt:
	goimports -w server

genrsa:
	cd server && \
	ssh-keygen -t rsa -b 2048 -m PEM -f rsa.key && \
	openssl rsa -in ./server/rsa.key -pubout -outform PEM -out ./server/rsa_pub.pem