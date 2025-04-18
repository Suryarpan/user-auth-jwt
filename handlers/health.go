package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Suryarpan/user-auth-jwt/utils"
)

type healthResp struct {
	Status string         `json:"status"`
	Config map[string]any `json:"config,omitempty"`
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	config := utils.NewConf()
	resp := healthResp{
		Status: "ok",
		Config: map[string]any{
			"debug":     config.Debug,
			"log_level": config.LogLevel,
		},
	}
	err := utils.Encode(w, http.StatusOK, resp)
	if err != nil {
		slog.Error("could not marshal", "error", err)
	}
}

func AttachHealthHandler(mux *http.ServeMux) {
	mux.HandleFunc("/health", healthHandler)
}
