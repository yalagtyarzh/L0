package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/internal/driver"
	"github.com/yalagtyarzh/L0/internal/handlers"
	"github.com/yalagtyarzh/L0/internal/msgbroker"
	"github.com/yalagtyarzh/L0/internal/render"
	"github.com/yalagtyarzh/L0/internal/repocache"
	"github.com/yalagtyarzh/L0/internal/repository"
	"github.com/yalagtyarzh/L0/internal/repository/dbrepo"
	"github.com/yalagtyarzh/L0/pkg/logging"
)

var app config.AppConfig
var cfg *config.Config
var logger *logging.Logger

func main() {
	log.Println("reading environment")
	cfg = config.GetConfig()

	log.Println("logger initializing")
	logger = logging.InitLogger(cfg.InProduction, cfg.AppConfig.LogLevel)
	defer logger.Sync()

	logger.Info("mux initializing")
	router := Router()

	logger.Info("connecting to database")
	db, err := connectDB()
	if err != nil {
		logger.Fatal(fmt.Sprintf("database error: %s", err))
	}
	defer db.SQL.Close()
	repo := dbrepo.NewPostgresRepo(db.SQL)

	logger.Info("recovering cache")
	cache := repocache.NewCache()
	err = cache.Recover(repo)
	if err != nil {
		logger.Fatal(fmt.Sprintf("error in recovering cache: %s", err))
	}

	logger.Info("creating template cache")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		logger.Fatal(fmt.Sprintf("cannot create template repocache: %s", err))
	}

	logger.Info("starting listening STAN")
	stan := msgbroker.NewSTAN(&cfg.STAN, cache, logger, repo)
	stan.SendMessages()

	setApp(tc)
	shareApp(repo, cache)

	logger.Info("Starting application")
	srv := &http.Server{
		Addr:    cfg.Listen.Port,
		Handler: router,
	}

	err = srv.ListenAndServe()
	logger.Fatal(err.Error())
}

func connectDB() (*driver.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s", cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port, cfg.PostgreSQL.Database,
	)

	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		return nil, err
	}

	logger.Info("connected!")

	return db, nil
}

func setApp(tc map[string]*template.Template) {
	app.InProduction = cfg.InProduction
	app.UseCache = cfg.UseCache
	app.TemplateCache = tc
}

func shareApp(repo repository.DatabaseRepo, cache *repocache.Cache) {
	r := handlers.NewRepo(&app, repo, cache)
	render.NewRenderer(&app)
	handlers.NewHandler(r)
}
