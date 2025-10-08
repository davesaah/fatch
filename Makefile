sqlc:
	@sqlc generate --file database/sqlc.yml

run:
	@go run cmd/server/main.go

dev:
	@air
