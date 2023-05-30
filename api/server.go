package api

import (
	"context"
	"log"
	"net/http"
	"simple-bank-system/db/services"

	"github.com/julienschmidt/httprouter"
)

// "Server" will serves all htpp request.
// services.Store to interact with db when processing API request from clients.

type Server struct {
	store *services.Store
	ctx   context.Context
}

// This func will create new "Server" instance, and setup all HTTP API routes for services on that server

func NewServer(store *services.Store, ctx context.Context) {
	address := ":8080"
	server := &Server{
		store: store,
		ctx:   ctx,
	}
	router := httprouter.New()

	// "createAccount" is made to be a method of the server, so it get access to the "store" object
	// in order to save new account ro the database
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	log.Println("Listening on localhost:", address)
	log.Fatal(http.ListenAndServe(":8080", router))
}
