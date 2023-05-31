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
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannpt connect to db: ", err)
	}

	dbpool, err := pgxpool.Connect(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("Failed connect to db: ", err)
	}
	defer dbpool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	store := services.NewStore(dbpool)
	api.NewServer(store, ctx, config.ServerAddress)
}
