package main

import (
	"database/sql"
	"github.com/isaya1910/zhasa-news/api"
	db "github.com/isaya1910/zhasa-news/db/sqlc"
	"github.com/isaya1910/zhasa-news/util"
	_ "github.com/lib/pq"
	"google.golang.org/api/option"
	"log"
	"os"
	"path/filepath"
)

var pathSearch string

func search() {
	err := filepath.Walk("/app", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == "serviceAccount.json" {
			pathSearch = path
			log.Println("Found file:", path)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
}
func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	search()
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db", err)
	}
	store := db.NewStore(conn)

	configPath := pathSearch
	log.Println(configPath)

	opt := option.WithCredentialsFile(configPath)

	server := api.NewServer(opt, store, api.UserExternalRepository{})

	err = server.Start(config.ServerAddress)

	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
