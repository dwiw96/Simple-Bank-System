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
	router.POST("/account", server.createAccount)
	router.POST("/account/login", server.loginAccount)

	// Add middleware auth to handler
	router.POST("/wallet", authMiddleware(server.tokenMaker, server.createWallet))
	router.GET("/wallet/:number", authMiddleware(server.tokenMaker, server.getWallet))
	router.GET("/wallet", authMiddleware(server.tokenMaker, server.listWallets))
	router.PUT("/wallet/update/:number", authMiddleware(server.tokenMaker, server.updateWallet))
	router.PUT("/wallet/updateInfo/:number", authMiddleware(server.tokenMaker, server.updateWalletInfo))
	router.DELETE("/wallet/delete/:number", authMiddleware(server.tokenMaker, server.deleteWallet))

	router.POST("/transfer", authMiddleware(server.tokenMaker, server.createTransfer))
	//router.GET("/transfer/:number", authMiddleware(server.tokenMaker, server.createTransfer))
	//router.GET("/transfer/list/:number", authMiddleware(server.tokenMaker, server.listTransfer))

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
