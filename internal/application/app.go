package application

import (
	"fmt"

	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/internal/driver"
	"github.com/yalagtyarzh/L0/internal/middleware"
	"github.com/yalagtyarzh/L0/internal/mux"
	"github.com/yalagtyarzh/L0/internal/render"
	"github.com/yalagtyarzh/L0/pkg/logging"
)

type Application struct {
	config.AppConfig
}

func NewApp(cfg *config.Config, logger *logging.Logger) (*Application, error) {
	logger.Info("mux initializing")
	router := mux.Router()

	logger.Info("connecting to database")
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s", cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Host,
		cfg.PostgreSQL.Port, cfg.PostgreSQL.Database,
	)

	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info("connected!")

	logger.Info("creating template cache")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		logger.Fatal(fmt.Sprintf("cannot create template cache: %s", err.Error()))
		return nil, err
	}

	appconfig := config.AppConfig{
		InProduction:  cfg.InProduction,
		UseCache:      cfg.UseCache,
		TemplateCache: tc,
		Logger:        logger,
		Router:        router,
		DB:            db,
	}

	return &Application{appconfig}, nil
}

// Run starting server for application TODO: Do http server xd
func (a *Application) Run() {
	shareApp(&a.AppConfig)
}

// shareApp shares appconfig for internal packages
func shareApp(appcfg *config.AppConfig) {
	render.NewRenderer(appcfg)
	middleware.NewMiddleware(appcfg)
}
