package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"simple-bank-system/db/services"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/julienschmidt/httprouter"
)

func translateError(err error, trans ut.Translator) (errs []string) {
	if err == nil {
		return nil
	}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(trans))
		errs = append(errs, translatedErr.Error())
	}
	return errs
}

type createAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required,currency"`
}

func (server *Server) createAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var req createAccountRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode input data", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	arg := services.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			http.Error(w, "Failed to encode error from validate func", (http.StatusInternalServerError))
			return
		}
		return
	}

	//translate
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	valErr := translateError(err, trans)
	if err != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(valErr)
		if err != nil {
			fmt.Println("Encode err: ", err)
		}
		return
	}

	account, err := server.store.CreateAccount(server.ctx, arg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, "failed to pass data into database", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

type getAccountRequest struct {
	ID int64
}

func (server *Server) getAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")
	var req getAccountRequest
	var err error

	req.ID, err = strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert url parameter to int", (http.StatusInternalServerError))
		return
	}

	account, err := server.store.GetAccount(server.ctx, req.ID)
	if err != nil {
		http.Error(w, "Can't get account", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(account)
}

type listAccountsRequest struct {
	PageID   int `URI:"page_id" validate:"required,min=1"`
	PageSize int `URI:"page_size" validate:"required,min=1,max=10"`
}

func (server *Server) listAccounts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	var req listAccountsRequest
	var err error

	req.PageID, err = strconv.Atoi(r.URL.Query().Get("page_id"))
	if err != nil {
		http.Error(w, "Failed convert page_id query to int", (http.StatusInternalServerError))
		return
	}
	req.PageSize, err = strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil {
		http.Error(w, "Failed convert page_size query to int", (http.StatusInternalServerError))
		return
	}

	arg := services.ListAccountParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	validate := validator.New()
	err = validate.Struct(req)

	//translate
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	valErr := translateError(err, trans)
	if valErr != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		json.NewEncoder(w).Encode(valErr)
		return
	}
	accounts, err := server.store.ListAccount(server.ctx, arg)
	if err != nil {
		http.Error(w, "Failed to get List", (http.StatusInternalServerError))
		return
	}
	if accounts == nil {
		http.Error(w, "There are no data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accounts)
}

type updateAccountRequest struct {
	ID      int64 `validate:"required,min=1"`
	Balance int64 `json:"balance" validate:"required"`
}

func (server *Server) updateAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req updateAccountRequest
	var err error

	req.ID, err = strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, "Failed convert data to int or data is nill", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	fmt.Println("id: ", req.ID)

	err = json.NewDecoder(r.Body).Decode(&req)
	fmt.Println("balance: ", req.Balance)
	if err != nil {
		http.Error(w, "Failed to get input data", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	arg := services.UpdateAccountParams{
		ID:      req.ID,
		Balance: req.Balance,
	}

	validate := validator.New()
	err = validate.Struct(req)

	//translate
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	valErr := translateError(err, trans)
	if valErr != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		json.NewEncoder(w).Encode(valErr)
		return
	}

	err = server.store.UpdateAccount(server.ctx, arg)
	if err != nil {
		http.Error(w, "Failed parsing data into database", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data modified")
}

func (server *Server) deleteAccount(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ID, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		http.Error(w, "Failed convert to int", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	validate := validator.New()
	err = validate.Var(ID, "required")

	//translate
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	valErr := translateError(err, trans)
	if valErr != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		json.NewEncoder(w).Encode(valErr)
		return
	}

	err = server.store.DeleteAccount(server.ctx, ID)
	if err != nil {
		http.Error(w, "Failed passing data to database", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data deleted!")
}
