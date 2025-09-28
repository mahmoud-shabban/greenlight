package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	logger *slog.Logger
	config config
}

func main() {

	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "dev", "server environment")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	mux.HandleFunc("/v1/healthcheck", app.healthChecker)

	logger.Info("starting server", slog.Any("addr", srv.Addr))

	err := srv.ListenAndServe()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
