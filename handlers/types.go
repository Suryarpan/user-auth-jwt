// type defs for request and responses
package handlers

import (
	"time"

	"github.com/Suryarpan/user-auth-jwt/database"
	"github.com/jackc/pgx/v5/pgtype"
)

// health handler
type healthResp struct {
	Status   string         `json:"status"`
	Config   map[string]any `json:"config,omitempty"`
	DbResult string         `json:"db_result"`
}

// public user details

type userPublic struct {
	UserId       pgtype.UUID `json:"user_id"`
	Username     string      `json:"username"`
	DisplayName  string      `json:"display_name"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	LastLoggedIn time.Time   `json:"last_logged_in"`
}

func newPublicUser(du database.User) userPublic {
	return userPublic{
		UserId:       du.UserId,
		Username:     du.Username,
		DisplayName:  du.DisplayName,
		CreatedAt:    du.CreatedAt,
		UpdatedAt:    du.UpdatedAt,
		LastLoggedIn: du.LastLoggedIn,
	}
}

// create user handler
type createUserReq struct {
	Username    string `json:"username" validate:"required,min=5,max=50"`
	DisplayName string `json:"display_name" validate:"required,min=5,max=150"`
	Password    string `json:"password" validate:"required,printascii,min=8"`
}
