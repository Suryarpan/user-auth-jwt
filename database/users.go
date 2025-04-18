package database

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

const createUser string = `INSERT INTO users (
    user_id, username, display_name, password, password_salt
) VALUES (
    gen_random_uuid(), $1, $2, $3, $4
) RETURNING pvt_id, user_id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in`

func CreateUser(r *http.Request, conn *pgxpool.Pool, arg CreateUserParams) (User, error) {
	row := conn.QueryRow(
		r.Context(),
		createUser,
		arg.Username,
		arg.DisplayName,
		arg.Password,
		arg.PasswordSalt,
	)
	var d User
	err := d.ScanRow(row)
	return d, err
}

const getUserByUUID string = `SELECT
	pvt_id,
	user_id,
	username,
	display_name,
	password,
	password_salt,
	created_at,
	updated_at,
	last_logged_in
FROM users
WHERE user_id = $1`

func GetUserByUUID(r *http.Request, conn *pgxpool.Pool, arg pgtype.UUID) (User, error) {
	row := conn.QueryRow(r.Context(), getUserByUUID, arg)
	var d User
	err := d.ScanRow(row)
	return d, err
}
