package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Suryarpan/user-auth-jwt/database"
	"github.com/Suryarpan/user-auth-jwt/middleware"
	"github.com/Suryarpan/user-auth-jwt/utils"
)

const (
	usedUser    string = "username is already taken"
	serverError string = "could not create new user"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	u := createUserReq{}
	llo := middleware.GetLLObject(r)
	err := utils.ValidateReq(w, r, llo.Validator, &u)
	if err != nil {
		return
	}
	// additional validation
	_, err = database.GetUserByUserName(r, llo.Conn, u.Username)
	if err == nil {
		slog.Error("username is already used", "name", u.Username, "display_name", u.DisplayName)
		utils.EncodeError(w, http.StatusConflict, usedUser)
		return
	}
	// password hash
	passwordSalt, err := utils.GenerateSalt()
	if err != nil {
		slog.Error("could not generate salt", "error", err)
		utils.EncodeError(w, http.StatusInsufficientStorage, serverError)
		return
	}
	password, err := utils.SaltyPassword([]byte(u.Password), passwordSalt)
	if err != nil {
		slog.Error("could not generate password", "error", err)
		utils.EncodeError(w, http.StatusInsufficientStorage, serverError)
		return
	}
	// store in db
	cu, err := database.CreateUser(r, llo.Conn, database.CreateUserParams{
		Username:     u.Username,
		DisplayName:  u.DisplayName,
		Password:     password,
		PasswordSalt: passwordSalt,
	})
	if err != nil {
		database.GenericErrorLogger(err, "could not create user")
		utils.EncodeError(w, http.StatusInsufficientStorage, serverError)
		return
	}
	utils.Encode(w, http.StatusCreated, newPublicUser(cu))
}

func AttachUserHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /user", createUserHandler)
}
