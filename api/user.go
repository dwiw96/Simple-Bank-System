package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"simple-bank-system/db/services"
	"simple-bank-system/util"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/julienschmidt/httprouter"
)

type createUserRequest struct {
	Username       string `json:"username" validate:"required"`
	HashedPassword string `json:"password" validate:"required,min=6"`
	FullName       string `json:"fullname" validate:"required"`
	Email          string `json:"email" validate:"required,email"`
}

type createUserResponse struct {
	Username         string
	FullName         string
	Email            string
	PasswordChangeAt time.Time
	CreatedAt        time.Time
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req createUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "failed to decode input data", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	arg := services.CreateUserParams{
		Username:       req.Username,
		HashedPassword: req.HashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
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

	user, err := server.store.CreateUser(server.ctx, arg)
	if err != nil {
		if err == util.ErrUser {
			http.Error(w, err.Error(), 403)
			return
		}
		http.Error(w, "failed to pass data into database", (http.StatusInternalServerError))
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	response := createUserResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
