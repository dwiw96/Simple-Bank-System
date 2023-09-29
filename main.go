package main

import (
	"context"
	"log"
	"time"

	"simple-bank-system/api"
	"simple-bank-system/db/services"
	"simple-bank-system/util"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	log.Println("--- Main()")
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Failed connect to db: ", err)
	}
	defer dbpool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Second)
	defer cancel()

	store := services.NewStore(dbpool)
	server, err := api.NewServer(store, ctx, config)
	if err != nil {
		log.Fatal("Can't create server, \nerr: ", err)
	}
	log.Println("--- (1) TestMain()")
	server.Start(config.ServerAddress)
	log.Println("--- (2) TestMain()")
}
