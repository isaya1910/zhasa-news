package main

import (
	"database/sql"
	"github.com/isaya1910/zhasa-news/api"

	db "github.com/isaya1910/zhasa-news/db/sqlc"
	_ "github.com/lib/pq"
	"log"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@localhost:5432/zhasa_news?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
