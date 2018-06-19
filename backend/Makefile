build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app
	docker build -t ${USER}:say .
	rm -rf app

run:
	docker run --rm -p 8080:8080 ${USER}:say