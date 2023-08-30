package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"simple-bank-system/token"

	"github.com/julienschmidt/httprouter"
)

// Function for CORS setting
func CORSMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if r.Method == http.MethodOptions {

			return
		}

		log.Println("-----(1) CORS handler")
		//w.Header().Set("Content-Type", "application/json; text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000, http://localhost:3000/register") // Change to your frontend URL
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			log.Println("OPTIONS")
			//http.Error(w, "No Content", http.StatusNoContent)
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000, http://localhost:3000/register") // Change to your frontend URL
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r, ps)
	}
}

func authMiddleware(tokenMaker token.Maker, next httprouter.Handle) httprouter.Handle {
	return (func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			http.Error(w, "Authorization header isn't provided", http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		authHeaderType := strings.ToLower(fields[0])
		if authHeaderType != "bearer" {
			http.Error(w, "Authorization header type is wrong", http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "token is wrong", http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), "authPayloadKey", payload)
		r = r.WithContext(ctx)

		// Call registered handler
		next(w, r, ps)
	})
}
