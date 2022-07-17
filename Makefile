build:
	GOOS=linux GOARCH=amd64 go build -o portdomain-linux

docker: build

	docker build -t portdomain:develop .

lint:
	golangci-lint  run