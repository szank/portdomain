build:
	GOOS=linux GOARCH=amd64 go build

docker:
	docker build -t portdomain:develop .

lint:
	golangci-lint  run