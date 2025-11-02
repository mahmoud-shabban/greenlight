package main

import (
	"flag"
	"os"

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

	limiter struct {
		rps     float64
		burst   int
		enabled bool
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
	flag.BoolVar(&cfg.limiter.enabled, "limiter", true, "Enable/Disable rate limitier (default true)")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 4, "Rate limiter maximum requests per second (default 4)")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 6, "Rate Limiter maximum burst ")

	flag.Parse()

	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.PrintError(err, nil)
		panic(1)
	}

	defer db.Close()

	logger.PrintInfo("DB connection pool stablished successfully.", nil)

	app := &Application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	err = app.serve()
	// srv := &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", app.config.port),
	// 	Handler:      app.routes(),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 30 * time.Second,
	// 	ErrorLog:     log.New(logger, "", 0),
	// }

	// logger.PrintInfo("starting server", map[string]string{
	// 	"addr": srv.Addr,
	// 	"env":  cfg.env,
	// })

	// err = srv.ListenAndServe()

	if err != nil {
		logger.PrintError(err, nil)
		panic(1)
	}

}
