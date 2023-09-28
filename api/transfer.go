package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"simple-bank-system/db/services"
	"simple-bank-system/token"
	"simple-bank-system/util"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/julienschmidt/httprouter"
)

type transferRequest struct {
	FromWalletNumber int64  `json:"from_wallet_number" validate:"required,min=1010000000,max=1019999999"`
	ToWalletNumber   int64  `json:"to_wallet_number" validate:"required,min=1010000000,max=1019999999"`
	Amount           int64  `json:"amount" validate:"required,gt=0"`
	Currency         string `json:"currency" validate:"required,currency"`
}

type transferResponse struct {
	FromWalletNumber int64
	ToWalletNumber   int64
	Amount           int64
	CreatedAt        time.Time
}

type entryResponse struct {
	WalletNumber int64
	Amount       int64
	CreatedAt    time.Time
}

type transferTxResponse struct {
	Transfer   transferResponse
	FromWallet walletResponse
	ToWallet   walletResponse
	FromEntry  entryResponse
	ToEntry    entryResponse
}

func newTransferTxResponse(tx *services.TransferTXResult) transferTxResponse {
	return transferTxResponse{
		Transfer: transferResponse{
			FromWalletNumber: tx.Transfer.FromWalletNumber,
			ToWalletNumber:   tx.Transfer.ToWalletNumber,
			Amount:           tx.Transfer.Amount,
			CreatedAt:        tx.Transfer.CreatedAt,
		},
		FromWallet: walletResponse{
			Name:         tx.FromWallet.Name,
			WalletNumber: tx.FromWallet.WalletNumber,
			Balance:      tx.FromWallet.Balance,
			Currency:     tx.FromWallet.Currency,
			CreatedAt:    tx.FromWallet.CreatedAt,
		},
		ToWallet: walletResponse{
			Name:         tx.ToWallet.Name,
			WalletNumber: tx.ToWallet.WalletNumber,
			Balance:      tx.ToWallet.Balance,
			Currency:     tx.ToWallet.Currency,
			CreatedAt:    tx.ToWallet.CreatedAt,
		},
		FromEntry: entryResponse{
			WalletNumber: tx.FromEntry.WalletNumber,
			Amount:       tx.FromEntry.Amount,
			CreatedAt:    tx.FromEntry.CreatedAt,
		},
		ToEntry: entryResponse{
			WalletNumber: tx.ToEntry.WalletNumber,
			Amount:       tx.ToEntry.Amount,
			CreatedAt:    tx.ToEntry.CreatedAt,
		},
	}
}

func (server *Server) createTransfer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req transferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Failed to decode json", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	err = validate.Struct(req)
	if err != nil {
		http.Error(w, "Format input data is wrong", (http.StatusBadRequest))
		err = json.NewEncoder(w).Encode(err.Error())
		if err != nil {
			http.Error(w, "Failed to encode error from validate func", (http.StatusInternalServerError))
			return
		}
		log.Println("---(1) api/create transfer, err:", err)
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
			log.Println("---(2) api/Encode err: ", err)
		}
		log.Println("--- (3) api/translate err:", err)
		return
	}

	authPayload := r.Context().Value("authPayloadKey").(*token.Payload)

	accountID, valid := server.validWallet(w, req.FromWalletNumber, req.Currency)
	if !valid {
		return
	}

	if *accountID != authPayload.AccountID {
		http.Error(w, "from_account isn't authorization", http.StatusUnauthorized)
		return
	}

	_, valid = server.validWallet(w, req.ToWalletNumber, req.Currency)
	if !valid {
		return
	}

	// Read wallet info to get both account & wallet ID
	wallet, err := server.store.GetWalletByNumber(server.ctx, req.FromWalletNumber)
	if err != nil {
		http.Error(w, "Failed to get the wallet", (http.StatusBadRequest))
	}

	arg := services.TransferTxParams{
		AccountID:        wallet.AccountID,
		WalletID:         wallet.ID,
		FromWalletNumber: req.FromWalletNumber,
		ToWalletNumber:   req.ToWalletNumber,
		Amount:           req.Amount,
	}

	accounts, err := server.store.TransferTx(server.ctx, arg)
	if err != nil {
		http.Error(w, "Failed to tranfer", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	accountsResponse := newTransferTxResponse(accounts)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accountsResponse)
}

/*
 * Compare input currency with the currenccy of the 'from_account' and 'to_account' to make sure
 * that they're all the same. This func check if an account with spesific ID exist, and itt's
 * currency matches with input currency.
 */

func (server *Server) validWallet(w http.ResponseWriter, number int64, currency string) (*int64, bool) {
	//log.Println("1")
	wallet, err := server.store.GetWalletByNumber(server.ctx, number)
	if err != nil {
		if err == util.ErrNotExist {
			http.Error(w, err.Error(), http.StatusNotFound)
			return nil, false
		}
		http.Error(w, "Failed to check your wallet", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return nil, false
	}

	if wallet.Currency != currency {
		err = fmt.Errorf("wallet [%d] currency mismatch: %s - %s", wallet.ID, wallet.Currency, currency)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return nil, false
	}

	return &wallet.AccountID, true
}

/*func(server *Server) listTransfer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error

	walletNumber, err := strconv.ParseInt(ps.ByName("number"), 10, 64)
	if err != nil {
		http.Error(w, "failed to convert url parameter to int", (http.StatusInternalServerError))
		return
	}

	server.store.ListTransfers()

	authPayload := r.Context().Value("authPayloadKey").(*token.Payload)
}*/
