package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"simple-bank-system/db/services"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/julienschmidt/httprouter"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" validate:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" validate:"required,min=1"`
	Amount        int64  `json:"amount" validate:"required,gt=0"`
	Currency      string `json:"currency" validate:"required,currency"`
}

func (server *Server) createTransfer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//log.Println("1")
	var req transferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to decode json", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	//log.Println("2")
	//log.Println("3")
	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			http.Error(w, "Failed to encode error from validate func", (http.StatusInternalServerError))
			return
		}
		//log.Println("4")
		return
	}
	//log.Println("5")
	//translate error from validator v10
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	valErr := translateError(err, trans)
	if valErr != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(valErr)
		if err != nil {
			log.Println("Encode err: ", err)
		}
		log.Println("6")
		return
	}

	if !server.validAccount(w, req.FromAccountID, req.Currency) {
		return
	}
	if !server.validAccount(w, req.ToAccountID, req.Currency) {
		return
	}

	arg := services.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	accounts, err := server.store.TransferTx(server.ctx, arg)
	if err != nil {
		http.Error(w, "Failed to tranfer", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

/*
 * Compare input currency with the currenccy of the 'from_account' and 'to_account' to make sure
 * that they're all the same. This func check if an account with spesific ID exist, and itt's
 * currency matches with input currency.
 */

func (server *Server) validAccount(w http.ResponseWriter, id int64, currency string) bool {
	//log.Println("1")
	account, err := server.store.GetAccount(server.ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "There are no data with your account ID", http.StatusNotFound)
			json.NewEncoder(w).Encode(err.Error())
			//w.Write([]byte("0"))
			return false
		}
		http.Error(w, "Failed to check your account", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return false
	}

	if account.Currency != currency {
		err = fmt.Errorf("account [%d] currency mismatch: %s - %s", account.ID, account.Currency, currency)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return false
	}

	return true
}
