package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/mahmoud-shabban/greenlight/internal/data"
)

var version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type Application struct {
	logger *slog.Logger
	config config
	models data.Models
}

func main() {
	// dsn := "postgres://greenlight:pass@127.0.0.1/greenlight?sslmode=disable"

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "dev", "server environment")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("GREENLIGHT_DB_DSN"), "postgres database connection string")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Error(err.Error())
		panic(1)
	}

	defer db.Close()

	app := &Application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info("starting server", slog.Any("addr", srv.Addr))

	err = srv.ListenAndServe()

	if err != nil {
		logger.Error(err.Error())
		panic(1)
	}

}
