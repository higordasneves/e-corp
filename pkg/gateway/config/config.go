package config

import (
	"fmt"
	"log"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Auth AuthConfig
	DB   DatabaseConfig
	HTTP HTTP
	MQ   RabbitMQConfig
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

type HTTP struct {
	Address string `env:"HTTP_ADDR" env-default:"0.0.0.0"`
	Port    string `env:"HTTP_PORT" env-default:"8080"`
}

type RabbitMQConfig struct {
	Host     string `env:"RABBITMQ_HOST" env-default:"localhost"`
	User     string `env:"RABBITMQ_USER" env-default:"guest"`
	Password string `env:"RABBITMQ_PASSWORD" env-default:"guest"`
	Port     string `env:"RABBITMQ_PORT" env-default:"5672"`
	Exchange string `env:"RABBITMQ_EXCHANGE" env-default:"ecorp"`
	Queue    string `env:"RABBITMQ_QUEUE" env-default:"ecorp.stream.accountCreation"`
	Bind     string `env:"RABBITMQ_BIND" env-default:"ecorp.accountCreation"`
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

func (mqCfg RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%v:%v@%v:%v", mqCfg.User, mqCfg.Password, mqCfg.Host, mqCfg.Port)
}
