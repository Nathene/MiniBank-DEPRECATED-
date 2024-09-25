build:
	@go build -o bin/gobank

run: build
	@./bin/gobank

test:
	@go test -v ./...

seed:
	@ go run scripts/seeds/createUsers/main.go

rmtable:
	@ go run scripts/seeds/dropTable/main.go