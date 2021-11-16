package main

import (
	"database/sql"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/http/router"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	dbCfg := config.DatabaseConfig{}
	dbCfg.LoadEnv()
	dbDNS := dbCfg.DNS()

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

	r := router.GetHTTPHandler(db)
	log.Fatal(http.ListenAndServe(":8080", r))

}
