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
	"github.com/rs/cors"
)

// "Server" will serves all htpp request.
// services.Store to interact with db when processing API request from clients.

type Server struct {
	store      *services.Store
	ctx        context.Context
	router     *httprouter.Router
	handler    http.Handler
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

	router, handler := server.setupRouter()
	server = &Server{
		router:  router,
		handler: handler,
	}
	return server, nil
}

// func (server *Server) setupRouter() *httprouter.Router {
func (server *Server) setupRouter() (*httprouter.Router, http.Handler) {
	router := httprouter.New()

	// "createAccount" is made to be a method of the server, so it get access to the "store" object
	// in order to save new account ro the database
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	// Add middleware auth to handler
	router.POST("/accounts", authMiddleware(server.tokenMaker, server.createAccount))
	router.GET("/accounts/:id", authMiddleware(server.tokenMaker, server.getAccount))
	router.GET("/accounts", authMiddleware(server.tokenMaker, server.listAccounts))
	router.PUT("/accounts/:id", authMiddleware(server.tokenMaker, server.updateAccount))
	router.DELETE("/accounts/:id", authMiddleware(server.tokenMaker, server.deleteAccount))

	router.POST("/transfers", authMiddleware(server.tokenMaker, server.createTransfer))

	handler := cors.Default().Handler(router)

	return router, handler
}

func (server *Server) Start(address string) {
	log.Println("Listening on localhost:", address)
	log.Fatal(http.ListenAndServe(address, server.handler))
}

func CreatePasetoMaker(key []byte) (token.Maker, error) {
	maker, err := token.NewPasetoMaker(key)
	if err != nil {
		return nil, err
	}

	return maker, err
}
