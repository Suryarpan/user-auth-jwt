package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Suryarpan/user-auth-jwt/database"
	"github.com/Suryarpan/user-auth-jwt/middleware"
	"github.com/Suryarpan/user-auth-jwt/utils"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	config := utils.NewConf()
	llo := middleware.GetLLObject(r)
	// Query DB
	res, err := database.HealthCheck(r, llo.PgConn, database.HealthChecks{Data: "hheelllloooo"})
	if err != nil {
		slog.Error("could not query database", "error", err)
		utils.EncodeError(w, http.StatusInternalServerError, "could not get data")
		return
	}
	// create result
	resp := healthResp{
		Status: "ok",
		Config: map[string]any{
			"debug":     config.Debug,
			"log_level": config.LogLevel,
		},
		DbResult: res.Data,
	}
	err = utils.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func AttachHealthHandler(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)
}
