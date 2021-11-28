package postgres

import (
	"context"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/sirupsen/logrus"
	"os"
	"testing"
	"time"
)

var dbTest *pgxpool.Pool

func TestMain(m *testing.M) {

	dbCfg := &config.DatabaseConfig{}
	dbCfg.LoadEnv()
	dbCfg.Host = "localhost"
	dbCfg.Name = "ecorp_test"

	log := logrus.New()

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
			"POSTGRES_PASSWORD=" + dbCfg.Password,
			"POSTGRES_USER=" + dbCfg.User,
			"POSTGRES_DB=" + dbCfg.Name,
			"listen_addresses = '*'",
		},
	})

	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	_ = resource.Expire(90) // Tell docker to hard kill the container in 90 seconds

	dbCfg.Port = resource.GetPort("5432/tcp")
	dbDNS := dbCfg.DNS()
	log.Info("Connecting to database on url: ", dbDNS)

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 90 * time.Second
	if err = pool.Retry(func() error {
		dbTest, err = pgxpool.Connect(context.Background(), dbDNS)
		if err != nil {
			return err
		}
		err = dbTest.Ping(context.Background())
		return err
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	migrationPath := "migrations"
	err = Migration(migrationPath, dbTest, log)
	if err != nil {
		log.WithError(err).Fatal(config.ErrMigrateDB)
	}

	defer dbTest.Close()
	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
