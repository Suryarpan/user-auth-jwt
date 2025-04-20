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

type UserPublic struct {
	UserId       pgtype.UUID `json:"user_id"`
	Username     string      `json:"username"`
	DisplayName  string      `json:"display_name"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	LastLoggedIn time.Time   `json:"last_logged_in"`
}

func (u *UserPublic) FromUser(du database.User) {
	u.UserId = du.UserId
	u.Username = du.Username
	u.DisplayName = du.DisplayName
	u.CreatedAt = du.CreatedAt
	u.UpdatedAt = du.UpdatedAt
	u.LastLoggedIn = du.LastLoggedIn
}

// create user handler
type createUserReq struct {
	Username    string `json:"user_name"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}
