package handlers

import (
	"net/http"

	"github.com/yalagtyarzh/L0/models"
	"github.com/yalagtyarzh/L0/sub/internal/config"
	"github.com/yalagtyarzh/L0/sub/internal/driver"
	"github.com/yalagtyarzh/L0/sub/internal/render"
	"github.com/yalagtyarzh/L0/sub/internal/repository"
	"github.com/yalagtyarzh/L0/sub/internal/repository/dbrepo"
	"github.com/yalagtyarzh/L0/sub/internal/templatedata"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App   *config.AppConfig
	DB    repository.DatabaseRepo
	Cache map[string]models.Order
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) (*Repository, error) {
	d := dbrepo.NewPostgresRepo(db.SQL)

	orders, err := d.GetOrders()
	if err != nil {
		return nil, err
	}

	cache := map[string]models.Order{}
	for _, order := range orders {
		cache[order.OrderUID] = order
	}

	return &Repository{
		App:   a,
		DB:    d,
		Cache: cache,
	}, nil
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Index is the main page handler
func Index(w http.ResponseWriter, r *http.Request) {
	render.Template(
		w, r, "index.page.gohtml", &templatedata.TemplateData{},
	)
}
