.SILENT:

build:
	go mod verify
	go mod tidy
	go build --race -o .bin/app cmd/app/main.go

run: build
	./.bin/app
