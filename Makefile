# Makefile
.PHONY: migrate-up migrate-down sqlc

DB_URL=postgresql://viet:123@localhost:5432/mpc?sslmode=disable

migrate-create:
	@read -p "Enter migration name: " name; \
	goose -dir internal/db/migrations create $$name sql

migrate-up:
	goose -dir internal/db/migrations postgres "$(DB_URL)" up

migrate-down:
	goose -dir internal/db/migrations postgres "$(DB_URL)" down

sqlc:
	sqlc generate

run:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go