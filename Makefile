migrateup:
	migrate -path internal/db/migration -database "postgresql://dev:secret@localhost:5432/main_db?sslmode=disable" -verbose up
migratedown:
	migrate -path internal/db/migration -database "postgresql://dev:secret@localhost:5432/main_db?sslmode=disable" -verbose down
sqlc:
	go tool sqlc generate
test:
	go test -v -cover -race ./...
audit:
	go vet ./...
	go tool staticcheck ./...
	go tool govulncheck

.PHONY: migrateup migratedown sqlc test audit
