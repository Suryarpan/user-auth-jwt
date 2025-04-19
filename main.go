package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/Suryarpan/user-auth-jwt/handlers"
	"github.com/Suryarpan/user-auth-jwt/middleware"
	"github.com/Suryarpan/user-auth-jwt/utils"
)

func setupLogger(config *utils.ConfigType) {
	logConf := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.LogLevel,
	})
	logger := slog.New(logConf)
	slog.SetDefault(logger)
}

func lookSignal(srv *http.Server, connChan chan any) {
	sigClose := make(chan os.Signal, 1)
	signal.Notify(sigClose, os.Interrupt)
	<-sigClose
	err := srv.Shutdown(context.Background())
	if err != nil {
		slog.Error("server shutdown", "error", err)
	}
	close(connChan)
}

func main() {
	config := utils.NewConf()
	// basic setup
	setupLogger(config)
	middleware.LLOSetup()
	defer middleware.LLOClose()
	// mux setup
	baseMux := http.NewServeMux()
	// attach handlers
	handlers.AttachHealthHandler(baseMux)
	// api versioning
	ms := middleware.ChainMiddleware(
		middleware.ReqLogger,
		middleware.LLOMiddleware,
	)
	v1router := http.NewServeMux()
	v1router.Handle("/api/v1/", ms(http.StripPrefix("/api/v1", baseMux)))

	// server setup
	var srv http.Server = http.Server{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), config.LogLevel),
		Handler:  v1router,
	}
	idleConnsChan := make(chan any)
	go lookSignal(&srv, idleConnsChan)

	slog.Info(
		"starting server",
		"host", config.Host,
		"port", config.Port,
		"debug", config.Debug,
		"log_level", config.LogLevel,
	)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("could not start the server", "error", err)
		os.Exit(1)
	}
	<-idleConnsChan
}
