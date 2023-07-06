package api

import (
	"encoding/json"
	"fmt"
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

type createUserRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type userResponse struct {
	Username         string
	FullName         string
	Email            string
	PasswordChangeAt time.Time
	CreatedAt        time.Time
}

// Create func to pass user data in separated func from 'createUser()'
// because it has pasword data
func newUserResponse(user *pkg.User) userResponse {
	return userResponse{
		Username:         user.Username,
		FullName:         user.FullName,
		Email:            user.Email,
		PasswordChangeAt: user.PasswordChangeAt,
		CreatedAt:        user.CreatedAt,
	}
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
		HashedPassword: req.Password,
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
	response := newUserResponse(user)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

type loginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type loginResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Println("(1) Pass")
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("----------- (1)")
		http.Error(w, "Failed to decode login data", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	fmt.Println("(2) Pass")
	fmt.Println("request:")
	fmt.Println(req)
	user, err := server.store.GetUser(server.ctx, req.Username)
	if err != nil {
		fmt.Println("----------- (2)")
		http.Error(w, "user with input username doesn't exist", (http.StatusBadRequest))
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	fmt.Println("(3) Pass")
	fmt.Println("Get User: ")
	fmt.Println(user)
	err = util.VerifyPassword(req.Password, user.HashedPassword)
	fmt.Println("(3) Check")
	if err != nil {
		fmt.Println("----------- (3)")
		http.Error(w, "Your password is wrong", http.StatusBadRequest)
		json.NewEncoder(w).Encode(err.Error())
		return
	}
	fmt.Println("(4) Pass")
	token, err := server.tokenMaker.CreateToken(req.Username, server.duration)
	if err != nil {
		fmt.Println("----------- (4)")
		http.Error(w, "Failed to create encryted token", http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err.Error())
		return
	}

	fmt.Println("(5) Pass")
	response := loginResponse{
		AccessToken: token,
		User:        newUserResponse(user),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
