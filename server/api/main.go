package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/jsonlog"
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
	logger *jsonlog.Logger
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

	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.PrintError(err, nil)
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
		ErrorLog:     log.New(logger, "", 0),
	}

	logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  cfg.env,
	})

	err = srv.ListenAndServe()

	if err != nil {
		logger.PrintError(err, nil)
		panic(1)
	}

}
