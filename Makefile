.SILENT:

build:
	go mod verify
	go mod tidy
	go build --race -o .bin/app cmd/app/main.go

run: build
	./.bin/app

force:
	migrate -path migrations -database "postgres://rat:@localhost:5432/time_tracker?sslmode=disable" force 1
