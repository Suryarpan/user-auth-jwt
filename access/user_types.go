package access

import (
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	UserPublic struct {
		PvtId        int
		Username     string
		DisplayName  string
		Password     []byte
		PasswordSalt []byte
	}

	User struct {
		UserPublic
		UserId       pgtype.UUID
		CreatedAt    time.Time
		UpdatedAt    time.Time
		LastLoggedIn time.Time
	}

	RedisToken struct {
		UserId       string    `redis:"user_id"`
		UserPvtId    int       `redis:"user_pvt_id"`
		Expiry       time.Time `redis:"-"`
		RefreshToken string    `redis:"-"`
	}
)

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
