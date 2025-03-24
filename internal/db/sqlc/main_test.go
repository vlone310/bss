package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var testStore *Store
var testDB *pgxpool.Pool
var container testcontainers.Container

var dbName = "test_db"
var dbUser = "test"
var dbPassword = "test"

func TestMain(m *testing.M) {
	ctx := context.Background()
	pg, err := postgres.Run(ctx,
		"postgres:17-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)

	if err != nil {
		log.Fatalf("can not start postgres container %v", err)
	}

	container = pg

	host, err := pg.Host(ctx)
	if err != nil {
		log.Fatalf("can not get container host %v", err)
	}

	port, err := pg.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("can not get container port %v", err)
	}

	dbSource := fmt.Sprintf("postgresql://test:test@%s:%s/test_db?sslmode=disable", host, port.Port())

	config, err := pgxpool.ParseConfig(dbSource)
	if err != nil {
		log.Fatalf("can not parse config %v", err)
	}

	config.MaxConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	connPool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatalf("can not connect to db %v", err)
	}

	testDB = connPool
	testStore = NewStore(connPool)

	runDBMigration(dbSource)

	code := m.Run()

	if err := cleanup(ctx); err != nil {
		log.Printf("could not cleanup resources: %v", err)
	}

	os.Exit(code)
}

func cleanup(ctx context.Context) error {
	// Close database connection
	if testDB != nil {
		testDB.Close()
	}

	// Terminate container
	if container != nil {
		return container.Terminate(ctx)
	}
	return nil
}

func runDBMigration(dbURL string) {
	migration, err := migrate.New(
		"file://../migration",
		dbURL,
	)
	if err != nil {
		log.Fatalf("cannot create new migrate instance: %v", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("failed to run migrate up: %v", err)
	}
}
