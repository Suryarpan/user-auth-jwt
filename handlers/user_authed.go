package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Suryarpan/user-auth-jwt/database"
	"github.com/Suryarpan/user-auth-jwt/middleware"
	"github.com/Suryarpan/user-auth-jwt/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	userNotFound  string = "user not available"
	userForbidden string = "cannot view other user details"
	someError     string = "could not process request"
	notLogout     string = "could not logout"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	du := middleware.GetUserData(r)
	user_id := pgtype.UUID{}
	err := user_id.Scan(r.PathValue("user_id"))
	if err != nil {
		return
	}
	llo := middleware.GetLLObject(r)
	ru, err := database.GetUserByUUID(r, llo.PgConn, user_id)
	// validation
	if err != nil {
		database.GenericErrorLogger(err, "user not found")
		utils.EncodeError(w, http.StatusNotFound, userNotFound)
		return
	} else if ru.PvtId != du.PvtId { // TODO add admin user check here
		slog.Warn("user requested other user details", "user_id", user_id, "requested_by", du.UserId, "user_name", du.Username)
		utils.EncodeError(w, http.StatusForbidden, userForbidden)
		return
	}
	err = utils.Encode(w, http.StatusOK, newPublicUser(du))
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	u := updateUserReq{}
	llo := middleware.GetLLObject(r)
	err := utils.ValidateReq(w, r, llo.Validator, &u)
	if err != nil {
		return
	}
	du := middleware.GetUserData(r)
	// check if changed
	u.DisplayName = utils.If(u.DisplayName != "", u.DisplayName, du.DisplayName)
	u.Username = utils.If(u.Username != "", u.Username, du.Username)
	// password hashing
	var password []byte
	if u.Password != "" {
		password, err = utils.SaltyPassword([]byte(u.Password), du.PasswordSalt)
		if err != nil {
			slog.Error("could not generate password", "error", err)
			utils.EncodeError(w, http.StatusInsufficientStorage, serverError)
			return
		}
	} else {
		password = du.Password
	}
	// store in db
	uu, err := database.UpdateUser(r, llo.PgConn, database.UpdateUserParams{
		Username:    u.Username,
		DisplayName: u.DisplayName,
		Password:    password,
		PvtId:       du.PvtId,
	})
	if err != nil {
		database.GenericErrorLogger(err, "could not create user")
		utils.EncodeError(w, http.StatusInternalServerError, serverError)
		return
	}
	err = utils.Encode(w, http.StatusOK, newPublicUser(uu))
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	du := middleware.GetUserData(r)
	llo := middleware.GetLLObject(r)
	rn, err := database.DeleteUser(r, llo.PgConn, du.PvtId)
	if err != nil {
		database.GenericErrorLogger(err, "could not delete user")
		utils.EncodeError(w, http.StatusInternalServerError, someError)
		return
	} else if rn == 0 {
		slog.Error("could not delete user due to not found", "requested_by", du.UserId, "user_name", du.Username)
		utils.EncodeError(w, http.StatusNotFound, userNotFound)
		return
	}
	// delete redis
	err = database.DeleteToken(r, llo.RedisConn, du.UserId)
	if err != nil {
		// DB is already done and the key will expire
		slog.Warn("could not delete redis tokens", "requested_by", du.UserId, "user_name", du.Username, "error", err)
	}
	err = utils.Encode(w, http.StatusOK, newPublicUser(du))
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	du := middleware.GetUserData(r)
	llo := middleware.GetLLObject(r)
	// delete redis
	err := database.DeleteToken(r, llo.RedisConn, du.UserId)
	if err != nil {
		slog.Error("could not delete redis tokens", "requested_by", du.UserId, "user_name", du.Username, "error", err)
		utils.EncodeError(w, http.StatusInternalServerError, notLogout)
		return
	}
	err = utils.Encode(w, http.StatusOK, newPublicUser(du))
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func AttachUserAuthedHandler(mux *http.ServeMux) {
	am := middleware.ChainMiddleware(
		middleware.Authentication,
	)
	mux.Handle("POST /logout", am(http.HandlerFunc(logoutUser)))
	mux.Handle("GET /user/{user_id}", am(http.HandlerFunc(getUser)))
	mux.Handle("PATCH /user", am(http.HandlerFunc(updateUser)))
	mux.Handle("DELETE /user", am(http.HandlerFunc(deleteUser)))
}
