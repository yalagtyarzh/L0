package application

import (
	"fmt"

	"github.com/go-chi/chi"

	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/internal/driver"
	"github.com/yalagtyarzh/L0/internal/mux"
	"github.com/yalagtyarzh/L0/pkg/logging"
)

type Application struct {
	cfg    *config.Config
	logger *logging.Logger
	router *chi.Mux
	db     *driver.DB
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

	return &Application{
		cfg:    cfg,
		logger: logger,
		router: router,
		db:     db,
	}, nil
}

func (a *Application) Run() {

}
