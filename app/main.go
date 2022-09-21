package main

import (
	"database/sql"
	"github.com/isaya1910/zhasa-news/api"
	"github.com/isaya1910/zhasa-news/util"

	db "github.com/isaya1910/zhasa-news/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load configuratios")
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store, api.UserExternalRepository{})

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
