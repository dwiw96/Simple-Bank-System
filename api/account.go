package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"simple-bank-system/db/pkg"
	"simple-bank-system/db/services"
	"simple-bank-system/util"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/julienschmidt/httprouter"
)

type createAccountRequest struct {
	Username    string        `json:"username" validate:"required,min=3"`
	Password    string        `json:"password" validate:"required,min=6"`
	FullName    string        `json:"fullname" validate:"required"`
	DateOfBirth string        `json:"date_of_birth" validate:"required"`
	Address     pkg.Addresses `json:"address" validate:"required"`
	Email       string        `json:"email" validate:"required,email"`
}

type accountResponse struct {
	AccountNumber    int64
	Username         string
	FullName         string
	DateOfBirth      string
	Address          addressResponse
	Email            string
	PasswordChangeAt time.Time
	CreatedAt        time.Time
}

type addressResponse struct {
	Provinces string `json:"province" validate:"required"`
	City      string `json:"city" validate:"required"`
	ZIP       int64  `json:"zip" validate:"required"`
	Street    string `json:"street" validate:"required"`
}

// Create func to pass user data in separated func from 'createUser()'
// because it has pasword data
func newaccountResponse(account *pkg.Account) accountResponse {
	return accountResponse{
		AccountNumber: account.AccountNumber,
		Username:      account.Username,
		FullName:      account.FullName,
		DateOfBirth:   account.DateOfBirth.Format(time.DateOnly),
		Address: addressResponse{
			Provinces: account.Address.Provinces,
			City:      account.Address.City,
			ZIP:       account.Address.ZIP,
			Street:    account.Address.Street,
		},
		Email:            account.Email,
		PasswordChangeAt: account.PasswordChangeAt,
		CreatedAt:        account.CreatedAt,
	}
}

func (server *Server) createAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//log.Println("----- Create new user handler")

	var req createAccountRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode input data", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	date, err := util.GetDOB(req.DateOfBirth)
	if err != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			http.Error(w, "Failed to encode error from validate func", (http.StatusUnprocessableEntity))
			return
		}
		return
	}

	arg := services.CreateAccountParams{
		Username:       req.Username,
		HashedPassword: req.Password,
		FullName:       req.FullName,
		DateOfBirth:    date,
		Address: pkg.Addresses{
			Provinces: req.Address.Provinces,
			City:      req.Address.City,
			ZIP:       req.Address.ZIP,
			Street:    req.Address.Street,
		},
		Email: req.Email,
	}

	err = validate.Struct(req)
	//log.Println("err:", err)
	if err != nil {
		http.Error(w, "Format input data is wrong", (http.StatusUnprocessableEntity))
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			http.Error(w, "Failed to encode error from validate func", (http.StatusUnprocessableEntity))
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
		http.Error(w, "Failed to translate", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(valErr)
		if err != nil {
			log.Println("Encode err: ", err)
			http.Error(w, err.Error(), (http.StatusBadRequest))
		}
		return
	}

	account, err := server.store.CreateAccount(server.ctx, arg)
	//log.Println("--- (done)Create account:", account)
	if err != nil {
		//log.Println("--- (err)CreateAccount")
		code := accErrHandling(err)
		http.Error(w, "failed to pass data into database", code)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	response := newaccountResponse(account)

	walletArg := services.CreateWalletParams{
		AccountID:    account.ID,
		Name:         "Primary Wallet",
		WalletNumber: account.AccountNumber,
		Currency:     "IDR",
		Balance:      1000000,
	}

	_, err = server.store.CreatePrimaryWallet(server.ctx, walletArg)
	if err != nil {
		if err == util.ErrUsernameEmpty {
			http.Error(w, err.Error(), 403)
			return
		}
		http.Error(w, "failed to pass data into database", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken string          `json:"access_token"`
	Account     accountResponse `json:"account"`
}

func (server *Server) loginAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("--- (1) login, err:", err)
		http.Error(w, "Failed to decode login data", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	account, err := server.store.GetAccount(server.ctx, req.Username)
	if err != nil {
		log.Println("--- (2) login, err:", err)
		http.Error(w, "account with input username doesn't exist", (http.StatusBadRequest))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = util.VerifyPassword(req.Password, account.HashedPassword)
	if err != nil {
		log.Println("--- (3) login, err:", err)
		http.Error(w, "Your password is wrong", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	token, err := server.tokenMaker.CreateToken(account.ID, server.duration)
	if err != nil {
		log.Println("--- (4) login, err:", err)
		http.Error(w, "Failed to create encryted token", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := loginResponse{
		AccessToken: token,
		Account:     newaccountResponse(account),
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func accErrHandling(err error) int {
	if err == util.ErrUsernameExists || err == util.ErrEmailExists {
		log.Println("--- (err)ErrHandling-1, err:", err)
		return http.StatusConflict
	}

	for i := range util.ErrReturn {
		if err == util.ErrReturn[i] {
			log.Println("--- (err)ErrHandling-2, err:", err)
			return http.StatusUnprocessableEntity
		}
	}
	return http.StatusBadRequest
}
