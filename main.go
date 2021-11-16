package main

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

func main() {
	dbCfg := config.DatabaseConfig{}
	dbCfg.LoadEnv()
	dbDNS := dbCfg.DNS()
	fmt.Println(uuid.NewString())

	fmt.Println("Accessing database", dbDNS)
	db, err := sql.Open(dbCfg.Driver, dbDNS)
	if err != nil {
		fmt.Println(err)
		log.Fatal(config.ErrConnectDB, err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println(err)
		log.Fatal(config.ErrConnectDB, err)
	}
}
