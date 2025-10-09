sqlc:
	@sqlc generate --file database/sqlc.yml

run:
	@go run cmd/server/main.go

dev:
	@air

docs:
	swag init -g cmd/server/main.go
