package main

import (
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mahmoud-shabban/greenlight/internal/data"
	"github.com/mahmoud-shabban/greenlight/internal/jsonlog"
	"github.com/mahmoud-shabban/greenlight/internal/mailer"
)

var (
	version   = "1.0.0"
	buildTime string
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}

	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	cors struct {
		trustedOrigins []string
	}
}

type Application struct {
	logger *jsonlog.Logger
	config config
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	// dsn := "postgres://greenlight:pass@127.0.0.1/greenlight?sslmode=disable"
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "server port to listen on")
	flag.StringVar(&cfg.env, "env", "dev", "server environment")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "postgres database connection string")
	// db config
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-cons", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-cons", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	// limiter settings
	flag.BoolVar(&cfg.limiter.enabled, "limiter", true, "Enable/Disable rate limitier (default true)")
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 4, "Rate limiter maximum requests per second (default 4)")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 6, "Rate Limiter maximum burst ")
	// smtp settings
	flag.StringVar(&cfg.smtp.host, "smtp-host", "live.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "api", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "56fd11b2ddb54923f2a81d1bf950c4d8", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <hello@demomailtrap.co>", "SMTP sender")
	// cors trusted origins
	flag.Func("cors-trusted-origins", "Truested CORS origins (space separated)", func(val string) error {

		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil

	})

	displayVersion := flag.Bool("version", false, "Display the version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\nBuild Time:\t%s\n", version, buildTime)
		os.Exit(0)
	}

	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintError(err, nil)
		panic(1)
	}

	defer db.Close()

	logger.PrintInfo("DB connection pool stablished successfully.", nil)

	expvar.NewString("version").Set(version)
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &Application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
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
