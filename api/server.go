package api

import (
	"context"
	"log"
	"net/http"
	"simple-bank-system/db/services"
	"simple-bank-system/token"
	"simple-bank-system/util"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
)

// "Server" will serves all htpp request.
// services.Store to interact with db when processing API request from clients.

type Server struct {
	store      *services.Store
	ctx        context.Context
	router     *httprouter.Router
	tokenMaker token.Maker
	duration   time.Duration
}

// This func will create new "Server" instance, and setup all HTTP API routes for services on that server
var validate *validator.Validate

func NewServer(store *services.Store, ctx context.Context, config util.Config) (*Server, error) {
	maker, err := CreatePasetoMaker([]byte(config.TokenSymmetricKey))
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:      store,
		ctx:        ctx,
		tokenMaker: maker,
		duration:   config.AccessTokenDuration,
	}

	validate = validator.New()
	validate.RegisterValidation("currency", validCurrency)

	router := server.setupRouter()
	server = &Server{
		router: router,
	}
	return server, nil
}

func (server *Server) setupRouter() *httprouter.Router {
	router := httprouter.New()

	// "createAccount" is made to be a method of the server, so it get access to the "store" object
	// in order to save new account ro the database
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfer)

	return router
}

func (server *Server) Start(address string) {
	log.Println("Listening on localhost:", address)
	log.Fatal(http.ListenAndServe(address, server.router))
}

func CreatePasetoMaker(key []byte) (token.Maker, error) {
	maker, err := token.NewPasetoMaker(key)
	if err != nil {
		return nil, err
	}

	return maker, err
}
