package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Suryarpan/user-auth-jwt/database"
	"github.com/Suryarpan/user-auth-jwt/middleware"
	"github.com/Suryarpan/user-auth-jwt/utils"
)

const (
	usedUser    string = "username is already taken"
	serverError string = "could not create new user"
	authFail    string = "username or password is incorrect"
	userFail    string = "could not find the user"
	loginFail   string = "could not login"
	refreshFail string = "could not refresh"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	u := createUserReq{}
	llo := middleware.GetLLObject(r)
	err := utils.ValidateReq(w, r, llo.Validator, &u)
	if err != nil {
		return
	}
	// additional validation
	_, err = database.GetUserByUserName(r, llo.PgConn, u.Username)
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
	cu, err := database.CreateUser(r, llo.PgConn, database.CreateUserParams{
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

func loginUserHandler(w http.ResponseWriter, r *http.Request) {
	u := loginUserReq{}
	llo := middleware.GetLLObject(r)
	err := utils.ValidateReq(w, r, llo.Validator, &u)
	if err != nil {
		return
	}
	// additional validation
	du, err := database.GetUserByUserName(r, llo.PgConn, u.Username)
	if err != nil {
		slog.Error("could not get username", "name", u.Username)
		utils.EncodeError(w, http.StatusUnauthorized, authFail)
		return
	}
	isValid := utils.IsPassword(du.Password, []byte(u.Password), du.PasswordSalt)
	if !isValid {
		slog.Error("incorrect password provided", "name", u.Username)
		utils.EncodeError(w, http.StatusUnauthorized, authFail)
		return
	}
	// tokens
	when := time.Now()
	accToken, accExpiry, err1 := middleware.UserToToken(du, when)
	rfsToken, rfsExpiry, err2 := middleware.RefreshToken(du, when)
	if err1 != nil || err2 != nil {
		slog.Error("could not generate tokens", "access_error", err1, "refresh_error", err2)
		utils.EncodeError(w, http.StatusInternalServerError, loginFail)
		return
	}
	err = database.StoreToken(r, llo.RedisConn, database.RedisToken{
		UserId:       du.UserId.String(),
		UserPvtId:    du.PvtId,
		Username:     du.Username,
		Expiry:       rfsExpiry,
		RefreshToken: rfsToken,
	})
	if err != nil {
		slog.Error("could not save refresh token", "error", err)
		utils.EncodeError(w, http.StatusInternalServerError, loginFail)
		return
	}
	resp := tokenResp{
		AccessToken:      accToken,
		AccessExpiresAt:  accExpiry,
		RefreshToken:     rfsToken,
		RefreshExpiresAt: rfsExpiry,
	}
	err = utils.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func refreshTokenHndler(w http.ResponseWriter, r *http.Request) {
	t := refreshTokenReq{}
	llo := middleware.GetLLObject(r)
	err := utils.ValidateReq(w, r, llo.Validator, &t)
	if err != nil {
		return
	}
	// check with DB
	v, err := database.GetToken(r, llo.RedisConn, t.RefreshToken)
	if err != nil {
		slog.Error("could not get refresh token", "error", err)
		utils.EncodeError(w, http.StatusUnauthorized, userFail)
		return
	}
	d, err := database.GetUserById(r, llo.PgConn, v.UserPvtId)
	if err != nil {
		slog.Error("could not get user data", "error", err)
		utils.EncodeError(w, http.StatusUnauthorized, userFail)
		return
	}
	// generate new tokens
	when := time.Now()
	accToken, accExpiry, err := middleware.UserToToken(d, when)
	if err != nil {
		slog.Error("could not generate tokens", "error", err)
		utils.EncodeError(w, http.StatusInternalServerError, refreshFail)
		return
	}
	resp := tokenResp{
		AccessToken:      accToken,
		AccessExpiresAt:  accExpiry,
		RefreshToken:     v.RefreshToken,
		RefreshExpiresAt: v.Expiry,
	}
	err = utils.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func AttachUserHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /user", createUserHandler)
	mux.HandleFunc("POST /login", loginUserHandler)
	mux.HandleFunc("POST /refresh", refreshTokenHndler)
}
