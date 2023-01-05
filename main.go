package main

import (
	"database/sql"
	"fmt"
	"github.com/isaya1910/zhasa-news/api"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"github.com/isaya1910/zhasa-news/util"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"log"
	"os"
	"path/filepath"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	store := db.NewStore(conn)

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	configPath := filepath.Join(wd, "serviceAccount.json")
	log.Println(configPath)

	opt := option.WithCredentialsFile(configPath)

	server := api.NewServer(opt, store, api.UserExternalRepository{})

	server.Opt = opt

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
