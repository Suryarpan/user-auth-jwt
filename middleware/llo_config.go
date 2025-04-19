// this is to handle long lived objects in the application
package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"time"

	"github.com/Suryarpan/user-auth-jwt/utils"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxLloConfigKey string
type LLObjects struct {
	Conn      *pgxpool.Pool
	Validator *validator.Validate
}

const lloConfigData ctxLloConfigKey = "USER_AUTH_LLO_CONFIGS"

var llos LLObjects

func setupDbCon() *pgxpool.Pool {
	config := utils.NewConf()
	dbConfig, err := pgxpool.ParseConfig(config.DbUrl)
	if err != nil {
		slog.Error("unable to parse db config", "error", err)
		os.Exit(1)
	}

	dbConfig.MaxConns = 10
	dbConfig.MinConns = 0
	dbConfig.MaxConnLifetimeJitter = time.Hour * 1
	dbConfig.MaxConnIdleTime = time.Minute * 5
	dbConfig.HealthCheckPeriod = time.Minute
	dbConfig.ConnConfig.ConnectTimeout = time.Second * 10

	connPool, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		slog.Error("could not establish connection", "error", err)
		os.Exit(1)
	}

	err = connPool.Ping(context.Background())
	if err != nil {
		slog.Error("cannot ping database", "error", err)
		os.Exit(1)
	}
	return connPool
}

func setupValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterTagNameFunc(
		func(field reflect.StructField) string {
			return field.Tag.Get("json")
		},
	)
	return validate
}

func LLOMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), lloConfigData, llos)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func LLOSetup() {
	llos = LLObjects{
		Conn:      setupDbCon(),
		Validator: setupValidator(),
	}
}

func LLOClose() {
	llos.Conn.Close()
}

func GetLLObject(r *http.Request) LLObjects {
	llo, ok := r.Context().Value(lloConfigData).(LLObjects)
	if !ok {
		slog.Error("LLO is corrupted or overwritten")
		os.Exit(1)
	}
	return llo
}
