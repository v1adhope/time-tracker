build:
	go mod verify
	go mod tidy
	go build --race -o .bin/app cmd/app/main.go
.SILENT: build

# see https://github.com/swaggo/swag
docs:
	swag init -g internal/controllers/v1/router.go
.PHONY: docs

run: docs build
	./.bin/app
.SILENT: run

force:
	migrate -path migrations -database "postgres://rat:@localhost:5432/time_tracker?sslmode=disable" force 1
