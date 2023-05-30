package main

import (
	"context"
	"log"
	"time"

	"simple-bank-system/api"
	"simple-bank-system/db/services"

	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	dbpool, err := pgxpool.Connect(context.Background(), "postgresql://db:secret@localhost:5432/bank?sslmode=disable")
	if err != nil {
		log.Fatal("Failed connect to db: ", err)
	}
	defer dbpool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Second)
	defer cancel()

	store := services.NewStore(dbpool)
	api.NewServer(store, ctx)
}
