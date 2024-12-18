package usecase_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"

	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres"
	"github.com/higordasneves/e-corp/pkg/gateway/postgres/dbpool"
	"github.com/higordasneves/e-corp/utils/logger"
)

var mainPool *pgxpool.Pool

func TestMain(m *testing.M) {
	cfg := config.DatabaseConfig{
		Driver:   "postgres",
		Host:     "localhost",
		Name:     fmt.Sprintf("db_%d", time.Now().UnixNano()),
		User:     "postgres",
		Password: "postgres",
		Port:     "5432",
		SSLMode:  "prefer",
	}

	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_PASSWORD=" + cfg.Password,
			"POSTGRES_USER=" + cfg.User,
			"POSTGRES_DB=" + cfg.Name,
			"listen_addresses = '*'",
		},
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	_ = resource.Expire(90) // Tell docker to hard kill the container in 90 seconds

	cfg.Port = resource.GetPort("5432/tcp")
	dbDNS := cfg.DNS()
	log.Println("Connecting to database on url: ", dbDNS)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 90 * time.Second
	if err = pool.Retry(func() error {
		mainPool, err = pgxpool.New(context.Background(), dbDNS)
		if err != nil {
			return fmt.Errorf("creating new pool: %w", err)
		}

		err = mainPool.Ping(context.Background())
		if err != nil {
			return fmt.Errorf("testing conection: %w", err)
		}

		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Run tests
	code := m.Run()

	mainPool.Close()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

// NewDB creates a new database named as a sanitized dbName. It returns a connection pool to this database.
// It must be called after StartDockerContainer.
func NewDB(t *testing.T) dbpool.Conn {
	t.Helper()

	ctx, err := logger.NewWithCtx(context.Background())
	require.NoError(t, err)

	if mainPool == nil {
		return dbpool.Conn{}
	}

	dbName := fmt.Sprintf("db_%d", time.Now().UnixNano())

	_, err = mainPool.Exec(context.Background(), "create database "+dbName)
	require.NoError(t, err)

	connString := strings.Replace(mainPool.Config().ConnString(), mainPool.Config().ConnConfig.Database, dbName, 1)
	pool, err := pgxpool.New(context.Background(), connString)
	require.NoError(t, err)

	err = pool.Ping(context.Background())
	require.NoError(t, err)

	migrationPath := "../../gateway/postgres/migrations"
	err = postgres.Migration(ctx, migrationPath, pool)
	require.NoError(t, err)

	t.Cleanup(func() {
		pool.Close()
		_, _ = pool.Exec(context.Background(), "drop database "+dbName)
	})

	return dbpool.NewConn(pool)
}
