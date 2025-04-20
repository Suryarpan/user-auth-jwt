package database

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

const healthCheckQuery string = `SELECT
$1 AS health`

func HealthCheck(r *http.Request, conn *pgxpool.Pool, arg HealthChecks) (HealthChecks, error) {
	row := conn.QueryRow(
		r.Context(),
		healthCheckQuery,
		arg.Data,
	)
	var d HealthChecks
	slog.Info("health check", "data", fmt.Sprintf("%v", row))
	err := row.Scan(&d.Data)
	return d, err
}
