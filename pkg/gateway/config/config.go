package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

var (
	ErrConnectDB = errors.New("failed to connect to database")
	ErrMigrateDB = errors.New("failed to migrate to database")
)

type Config struct {
	Auth AuthConfig
	DB   DatabaseConfig
}

type AuthConfig struct {
	Duration  time.Duration `env:"AUTH_DURATION" env-default:"60m"`
	SecretKey string        `env:"AUTH_KEY" env-default:"replace_with_your_key"`
}

type DatabaseConfig struct {
	Driver   string `env:"DB_DRIVER" env-default:"postgres"`
	Host     string `env:"DB_HOST" env-default:"fullstack-postgres"`
	Name     string `env:"DB_NAME" env-default:"ecorp"`
	User     string `env:"DB_USER" env-default:"developer"`
	Password string `env:"DB_PASSWORD" env-default:"123456"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	SSLMode  string `env:"DB_SSL_MODE" env-default:"disable"`
}

// LoadEnv loads environment variables into a DatabaseConfig struct.
func (config *Config) LoadEnv() {
	err := cleanenv.ReadEnv(config)
	if err != nil {
		log.Fatal("env file not found")
	}
}

// DNS returns database domain name servers
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
