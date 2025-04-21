package handlers

import (
	"net/http"

	"github.com/Suryarpan/user-auth-jwt/middleware"
	"github.com/jackc/pgx/v5/pgtype"
)

func getUser(w http.ResponseWriter, r *http.Request) {
	du := middleware.GetUserData(r)
	user_id := pgtype.UUID{}
	err := user_id.Scan(r.PathValue("user_id"))
	if err != nil {
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	du := middleware.GetUserData(r)
	user_id := pgtype.UUID{}
	err := user_id.Scan(r.PathValue("user_id"))
	if err != nil {
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	du := middleware.GetUserData(r)
	user_id := pgtype.UUID{}
	err := user_id.Scan(r.PathValue("user_id"))
	if err != nil {
		return
	}
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	// this will not work with current refresh model
	du := middleware.GetUserData(r)
}

func AttachUserAuthedHandler(mux *http.ServeMux) {
	am := middleware.ChainMiddleware(
		middleware.Authentication,
	)
	mux.Handle("POST /logout", am(http.HandlerFunc(logoutUser)))
	mux.Handle("GET /user/{user_id}", am(http.HandlerFunc(getUser)))
	mux.Handle("PATCH /user/{user_id}", am(http.HandlerFunc(updateUser)))
	mux.Handle("DELETE /user/{user_id}", am(http.HandlerFunc(deleteUser)))
}
