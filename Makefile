.PHONY: migrate server

migrate:
    go run migrations/migrate.go

server:
    go run cmd/server/main.go
