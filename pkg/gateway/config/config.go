package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

var (
	ErrConnectDB = errors.New("failed to connect to database")
	ErrMigrateDB = errors.New("failed to migrate to database")
)

type DatabaseConfig struct {
	Driver   string `env:"DB_DRIVER" env-default:"postgres"`
	Host     string `env:"DB_HOST" env-default:"fullstack-postgres"`
	Name     string `env:"DB_NAME" env-default:"ecorp"`
	User     string `env:"DB_USER" env-default:"developer"`
	Password string `env:"DB_PASSWORD" env-default:"123456"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

//LoadEnv loads environment variables into a DatabaseConfig struct
func (dbConfig *DatabaseConfig) LoadEnv() {
	err := cleanenv.ReadEnv(dbConfig)
	if err != nil {
		log.Fatal("Fail, .env file not found")
	}
}

//DNS returns database domain name servers
func (dbConfig *DatabaseConfig) DNS() (dns string) {
	dns = fmt.Sprintf("%v://%v:%v@%v:%v/%v?sslmode=%v",
		dbConfig.Driver,
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Name,
		dbConfig.SSLMode,
	)
	return
}
