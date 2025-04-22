package database

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// DB User
type User struct {
	PvtId        int
	UserId       pgtype.UUID
	Username     string
	DisplayName  string
	Password     []byte
	PasswordSalt []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	LastLoggedIn time.Time
}

func (u *User) ScanRow(row pgx.Row) error {
	err := row.Scan(
		&u.PvtId,
		&u.UserId,
		&u.Username,
		&u.DisplayName,
		&u.Password,
		&u.PasswordSalt,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoggedIn,
	)
	return err
}

// Query Param structs
type HealthChecks struct {
	Data string
}

type CreateUserParams struct {
	Username     string
	DisplayName  string
	Password     []byte
	PasswordSalt []byte
}

type UpdateUserParams struct {
	Username    string
	DisplayName string
	Password    []byte
	PvtId       int
}

// Redis Token structs

type RedisToken struct {
	UserId       string    `redis:"user_id"`
	UserPvtId    int       `redis:"user_pvt_id"`
	Expiry       time.Time `redis:"-"`
	RefreshToken string    `redis:"-"`
}
