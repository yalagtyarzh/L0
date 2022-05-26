package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/yalagtyarzh/L0/sub/internal/config"
	"github.com/yalagtyarzh/L0/sub/internal/driver"
	"github.com/yalagtyarzh/L0/sub/internal/handlers"
	"github.com/yalagtyarzh/L0/sub/internal/render"
	"github.com/yalagtyarzh/L0/sub/pkg/logging"
)

var app config.AppConfig
var cfg *config.Config
var logger *logging.Logger

func main() {
	log.Println("reading environment")
	cfg = config.GetConfig()

	log.Println("logging initializing")
	logger = logging.InitLogger(cfg.InProduction, cfg.AppConfig.LogLevel)
	defer logger.Sync()

	logger.Info("mux initializing")
	router := Router()

	logger.Info("connecting to database")
	db, err := connectDB()
	if err != nil {
		logger.Fatal(fmt.Sprintf("database error: %s", err.Error()))
	}
	defer db.SQL.Close()

	logger.Info("creating template cache")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		logger.Fatal(fmt.Sprintf("cannot create template cache: %s", err.Error()))
	}

	setApp(router, tc)
	shareApp(db)

	logger.Info("Starting application")
	srv := &http.Server{
		Addr:    cfg.Listen.Port,
		Handler: app.Router,
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

func setApp(r http.Handler, tc map[string]*template.Template) {
	app.InProduction = cfg.InProduction
	app.UseCache = cfg.UseCache
	app.Router = r
	app.TemplateCache = tc
}

func shareApp(db *driver.DB) {
	repo := handlers.NewRepo(&app, db)
	render.NewRenderer(&app)
	handlers.NewHandler(repo)
}
