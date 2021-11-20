package main

import (
	"context"
	"fmt"
	"github.com/higordasneves/e-corp/pkg/gateway/config"
	"github.com/higordasneves/e-corp/pkg/gateway/http/router"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func main() {
	dbCfg := config.DatabaseConfig{}
	dbCfg.LoadEnv()
	dbDNS := dbCfg.DNS()

	fmt.Println("Accessing database", dbDNS)
	ctxDB := context.Background()
	dbPool, err := pgxpool.Connect(ctxDB, dbDNS)
	defer dbPool.Close()

	if err != nil {
		fmt.Println(err)
		log.Fatal(config.ErrConnectDB, err)
	}

	if err = dbPool.Ping(ctxDB); err != nil {
		fmt.Println(err)
		log.Fatal(config.ErrConnectDB, err)
	}

	r := router.GetHTTPHandler(dbPool)
	log.Fatal(http.ListenAndServe(":8080", r))

}
