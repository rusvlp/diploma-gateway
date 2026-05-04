USER := "postgres_user"
PASSWORD := "superStrongPassword"
HOST := "database"
PORT := "5432"
DB_NAME := "diploma"

DB_DSN := "postgres://${USER}:${PASSWORD}@${HOST}:${PORT}/${DB_NAME}?sslmode=disable"
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}
migrate:
	$(MIGRATE) up
migrate-down:
	$(MIGRATE) down
migrate-force:
	$(MIGRATE) force ${VERSION}
run:
	go run cmd/main.go
swagger:
	swag init -g cmd/main.go -o internal/docs