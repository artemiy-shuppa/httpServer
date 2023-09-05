include .env

clean:
	rm -f ./build/apiServer

build: clean
	go build -o ./build/ ./cmd/apiServer/

run: build
	./build/apiServer

test:
	go test -v -race -timeout 5s ./...

migrateup:
	@migrate -path ./migrations/ -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/app" -verbose up 

migratedown:
	@migrate -path ./migrations/ -database "postgresql://${DB_USER}:${DB_PASS}@${DB_HOST}/app" -verbose down 

.PHONY: build run test migrateup migratedown
.DEFAULT_GOAL := run
