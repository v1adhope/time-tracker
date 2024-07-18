build:
	go mod verify
	go mod tidy
	go build -o .bin/app cmd/app/main.go
.SILENT: build

docs:
	swag init -g internal/controllers/v1/router.go
.PHONY: docs

test:
	go test ./... --race

testByName:
	go test --race -run ${name} ./...

run: docs test build
	./.bin/app

force:
	migrate -path migrations -database "postgres://rat:@localhost:5432/time_tracker?sslmode=disable" force 1
