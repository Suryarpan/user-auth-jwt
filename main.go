package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func setupLogger(config *ConfigType) {
	logConf := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     config.LogLevel,
	})
	logger := slog.New(logConf)
	slog.SetDefault(logger)
}

func lookSignal(srv *http.Server, connChan chan any) {
	sigClose := make(chan os.Signal, 1)
	signal.Notify(sigClose, syscall.SIGINT, syscall.SIGTERM)
	<-sigClose
	err := srv.Shutdown(context.Background())
	if err != nil {
		slog.Error("server shutdown", "error", err)
	}
	close(connChan)
}

func main() {
	config := NewConf()
	// basic setup
	setupLogger(config)
	// mux setup
	mux := http.NewServeMux()
	// server setup
	var srv http.Server = http.Server{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		ErrorLog: slog.NewLogLogger(slog.Default().Handler(), config.LogLevel),
		Handler:  mux,
	}
	idleConnsChan := make(chan any)
	go lookSignal(&srv, idleConnsChan)
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("could not start the server", "error", err)
		os.Exit(1)
	}
	<-idleConnsChan
}
