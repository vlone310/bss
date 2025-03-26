migrateup:
	migrate -path internal/db/migration -database "postgresql://dev:secret@localhost:5432/main_db?sslmode=disable" -verbose up
migratedown:
	migrate -path internal/db/migration -database "postgresql://dev:secret@localhost:5432/main_db?sslmode=disable" -verbose down
sqlc:
	go tool sqlc generate
mockgen:
	go tool mockgen -destination internal/db/mock/store.go -package mockdb github.com/vlone310/bss/internal/db/sqlc Store
test:
	go test -v -cover -race ./...
audit:
	go vet ./...
	go tool staticcheck ./...
	go tool govulncheck
run:
	go run cmd/http/main.go
build:
	go build -v -ldflags "-s -w" -o bin/main cmd/http/main.go

.PHONY: migrateup migratedown sqlc test audit run mockgen
