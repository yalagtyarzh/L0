package handlers

import (
	"net/http"

	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/internal/driver"
	"github.com/yalagtyarzh/L0/internal/models"
	"github.com/yalagtyarzh/L0/internal/render"
	"github.com/yalagtyarzh/L0/internal/repository"
	"github.com/yalagtyarzh/L0/internal/repository/dbrepo"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Index is the main page handler
func Index(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "index.page.gohtml", &models.TemplateData{})
}
