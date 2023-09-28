package api

import (
	"encoding/json"
	"fmt"
	"log"

	"net/http"
	"strconv"
	"time"

	"simple-bank-system/db/pkg"
	"simple-bank-system/db/services"
	"simple-bank-system/token"
	"simple-bank-system/util"

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

type createWalletRequest struct {
	Name     string `json:"name" validate:"required"`
	Currency string `json:"currency" validate:"required,currency"`
}

type walletResponse struct {
	Name         string
	WalletNumber int64
	Balance      int64
	Currency     string
	CreatedAt    time.Time
}

func newWalletResponse(wallet *pkg.Wallet) walletResponse {
	return walletResponse{
		Name:         wallet.Name,
		WalletNumber: wallet.WalletNumber,
		Balance:      wallet.Balance,
		Currency:     wallet.Currency,
		CreatedAt:    wallet.CreatedAt,
	}
}

func (server *Server) createWallet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req createWalletRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode input data", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// access authorization payload inside context
	// 'ctx.Value' return general interface, so we should cast it to token.payload object
	authPayload := r.Context().Value("authPayloadKey").(*token.Payload)

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
		http.Error(w, "Format input data is wrong", (http.StatusUnprocessableEntity))
		err = json.NewEncoder(w).Encode(valErr)
		if err != nil {
			http.Error(w, "Failed to encode translation", (http.StatusUnprocessableEntity))
		}
		return
	}

	accountNumber, err := server.store.GetAccountByID(server.ctx, authPayload.AccountID)
	if err != nil {
		http.Error(w, "failed to create wallet because of account number", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	multiplier := 1
	var sixDigitNumb int64
	for i := 0; i < 10; i++ {
		if i > 3 {
			temp := *accountNumber % 10
			sixDigitNumb += temp * int64(multiplier)
			multiplier *= 10
		}
		*accountNumber /= 10
	}

	arg := services.CreateWalletParams{
		AccountID:    authPayload.AccountID, // add authorization to create account handler
		WalletNumber: sixDigitNumb * 10000,
		Name:         req.Name,
		Currency:     req.Currency,
		Balance:      0,
	}

	wallet, err := server.store.CreateWallet(server.ctx, arg)
	if err != nil {
		switch err {
		case util.ErrAccUser, util.ErrDuplicate:
			http.Error(w, err.Error(), 403)
			return
		}
		http.Error(w, "failed to pass data into database", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	response := newWalletResponse(wallet)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type getWalletRequest struct {
	WalletNumber int64
}

func (server *Server) getWallet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req getWalletRequest
	var err error

	req.WalletNumber, err = strconv.ParseInt(ps.ByName("number"), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert url parameter to int", (http.StatusInternalServerError))
		return
	}

	wallet, err := server.store.GetWalletByNumber(server.ctx, req.WalletNumber)
	if err != nil {
		if err == util.ErrNotExist {
			http.Error(w, err.Error(), 403)
			return
		}
		http.Error(w, "Can't get wallet", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// access authorization payload inside context
	// 'ctx.Value' return general interface, so we should cast it to token.payload object
	authPayload := r.Context().Value("authPayloadKey").(*token.Payload)
	if wallet.AccountID != authPayload.AccountID {
		http.Error(w, "wallet doesn't belong to you", (http.StatusUnauthorized))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallet)
}

type listWalletsRequest struct {
	PageID   int `URI:"page_id" validate:"required,min=1"`
	PageSize int `URI:"page_size" validate:"required,min=1,max=10"`
}

func (server *Server) listWallets(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Content-Type", "application/json")

	var req listWalletsRequest
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

	authPayload := r.Context().Value("authPayloadKey").(*token.Payload)

	arg := services.ListWalletParams{
		AccountID: authPayload.AccountID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
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
	wallets, err := server.store.ListWallet(server.ctx, arg)
	if err != nil {
		if err == util.ErrNotExist {
			http.Error(w, err.Error(), 503)
			return
		}
		http.Error(w, "Failed to get List", (http.StatusInternalServerError))
		return
	}
	if wallets == nil {
		http.Error(w, "There are no data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(wallets)
}

type updateWalletRequest struct {
	WalletNumber int64 `validate:"required,min=1010000000,max=1019999999"`
	Balance      int64 `json:"balance" validate:"required"`
}

func (server *Server) updateWallet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req updateWalletRequest
	var err error

	req.WalletNumber, err = strconv.ParseInt(ps.ByName("number"), 10, 64)
	if err != nil {
		http.Error(w, "Failed convert data to int or data is nill", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//fmt.Println("number: ", req.WalletNumber)

	err = json.NewDecoder(r.Body).Decode(&req)
	//fmt.Println("balance: ", req.Balance)
	if err != nil {
		http.Error(w, "Failed to get input data", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	wallet, err := server.store.GetWalletByNumber(server.ctx, req.WalletNumber)
	if err != nil {
		http.Error(w, "Failed to get the wallet", (http.StatusBadRequest))
	}

	arg := services.UpdateWalletParams{
		ID:      wallet.ID,
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

	err = server.store.UpdateWallet(server.ctx, arg)
	if err != nil {
		if err == util.ErrUpdateFailed {
			http.Error(w, err.Error(), 403)
			return
		}
		http.Error(w, "Failed parsing data into database", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data modified")
}

type updateWalletInfoRequest struct {
	WalletNumber int64  `validate:"required,min=1010000000,max=1019999999"`
	Name         string `json:"name" validate:"required"`
	Currency     string `json:"currency" validate:"required"`
}

func (server *Server) updateWalletInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req updateWalletInfoRequest
	var err error

	req.WalletNumber, err = strconv.ParseInt(ps.ByName("number"), 10, 64)
	if err != nil {
		http.Error(w, "Failed convert data to int or data is nill", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	//fmt.Println("number: ", req.WalletNumber)

	err = json.NewDecoder(r.Body).Decode(&req)
	//fmt.Println("balance: ", req.Balance)
	if err != nil {
		http.Error(w, "Failed to get input data", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	arg := services.UpdateWalletInformationParams{
		WalletNumber: req.WalletNumber,
		Name:         req.Name,
		Currency:     req.Currency,
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

	err = server.store.UpdateWalletInformation(server.ctx, arg)
	if err != nil {
		if err == util.ErrUpdateFailed {
			http.Error(w, err.Error(), 403)
			return
		}
		http.Error(w, "Failed parsing data into database", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Data modified")
}

func (server *Server) deleteWallet(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	Number, err := strconv.ParseInt(ps.ByName("number"), 10, 64)
	if err != nil {
		http.Error(w, "Failed convert to int", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	validate := validator.New()
	err = validate.Var(Number, "required")

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

	// read wallet information to get wallet ID, currency and balance
	wallet, err := server.store.GetWalletByNumber(server.ctx, Number)
	if err != nil {
		http.Error(w, "Failed to get the wallet", (http.StatusBadRequest))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	// Before delete the wallet, check the wallet balance.
	// if the wallet balance > 0, we should transfer the balance first to primary wallet using account number.
	_, isOK := server.transferWalletBalance(w, *wallet)
	if isOK == false {
		log.Println("---(error)transfer wallet balance")
		// http.Error(w, "Failed to delete wallet because transfer balance from wallet that will be deleted", http.StatusBadRequest)
		// json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = server.store.DeleteWallet(server.ctx, wallet.ID)
	if err != nil {
		http.Error(w, "Failed delete wallet from database", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Wallet deleted!")
}

func (server *Server) transferWalletBalance(w http.ResponseWriter, wallet pkg.Wallet) (*services.TransferTXResult, bool) {
	if wallet.Currency != "IDR" {
		err := fmt.Errorf("wallet [%d] currency mismatch: %s - %s", wallet.ID, wallet.Currency, "IDR")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return nil, false
	}

	accountNumber, err := server.store.GetAccountByID(server.ctx, wallet.AccountID)
	if err != nil {
		//log.Printf("(output) account ID: %d <> %d (input) account ID\n", *accountID, wallet.AccountID)
		http.Error(w, "failed to get account number to transfer wallet balance", http.StatusUnauthorized)
		return nil, false
	}

	arg := services.TransferTxParams{
		AccountID:        wallet.AccountID,
		WalletID:         wallet.ID,
		FromWalletNumber: wallet.WalletNumber,
		ToWalletNumber:   *accountNumber,
		Amount:           wallet.Balance,
	}

	accounts, err := server.store.TransferTx(server.ctx, arg)
	if err != nil {
		http.Error(w, "Failed to tranfer", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return nil, false
	}
	return accounts, true
}
