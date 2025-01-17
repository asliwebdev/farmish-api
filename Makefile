swag:
	go run github.com/swaggo/swag/cmd/swag init -g ./internal/handlers/router.go -o ./docs

DB_URL=postgres://postgres:Aa2210057@localhost:5432/postgres?sslmode=disable

migrate-up:
	goose -dir ./migrations postgres $(DB_URL) up

migrate-down:
	goose -dir ./migrations postgres $(DB_URL) down

migrate-status:
	goose -dir ./migrations postgres $(DB_URL) status