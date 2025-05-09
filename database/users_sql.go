package database

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
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

const getUserByUUID string = `SELECT pvt_id, user_id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
FROM users
WHERE user_id = $1`

func GetUserByUUID(r *http.Request, conn *pgxpool.Pool, arg pgtype.UUID) (User, error) {
	row := conn.QueryRow(r.Context(), getUserByUUID, arg)
	var d User
	err := d.ScanRow(row)
	return d, err
}

const getUserByUsername string = `SELECT pvt_id, user_id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
FROM users
WHERE username = $1`

func GetUserByUserName(r *http.Request, conn *pgxpool.Pool, arg string) (User, error) {
	row := conn.QueryRow(r.Context(), getUserByUsername, arg)
	var d User
	err := d.ScanRow(row)
	return d, err
}

const getUserById string = `SELECT pvt_id, user_id, username, display_name, password, password_salt, created_at, updated_at, last_logged_in
FROM users
WHERE pvt_id = $1`

func GetUserById(r *http.Request, conn *pgxpool.Pool, arg int) (User, error) {
	row := conn.QueryRow(r.Context(), getUserById, arg)
	var d User
	err := d.ScanRow(row)
	return d, err
}

const updateUser string = `UPDATE users
SET username = $1, display_name = $2, password = $3
WHERE pvt_id = $4
RETURNING *`

func UpdateUser(r *http.Request, conn *pgxpool.Pool, arg UpdateUserParams) (User, error) {
	row := conn.QueryRow(
		r.Context(),
		updateUser,
		arg.Username,
		arg.DisplayName,
		arg.Password,
		arg.PvtId,
	)
	var d User
	err := d.ScanRow(row)
	return d, err
}

const deleteUser string = `DELETE FROM users
WHERE pvt_id = $1`

func DeleteUser(r *http.Request, conn *pgxpool.Pool, arg int) (int64, error) {
	row, err := conn.Exec(r.Context(), deleteUser, arg)
	return row.RowsAffected(), err
}

func GenericErrorLogger(err error, mssg string) {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		slog.Error(
			mssg,
			"error", pgErr.Message,
			"code", pgErr.Code,
			"constraint", pgErr.ConstraintName,
		)
	} else {
		slog.Error(mssg)
	}
}
