package handlers

import (
	"net/http"

	"github.com/yalagtyarzh/L0/internal/config"
	"github.com/yalagtyarzh/L0/internal/render"
	"github.com/yalagtyarzh/L0/internal/repocache"
	"github.com/yalagtyarzh/L0/internal/repository"
	"github.com/yalagtyarzh/L0/internal/templatedata"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App   *config.AppConfig
	DB    repository.DatabaseRepo
	Cache *repocache.Cache
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, repo repository.DatabaseRepo, c *repocache.Cache) *Repository {
	return &Repository{
		App:   a,
		DB:    repo,
		Cache: c,
	}
}

// NewHandler sets the repository for the handlers
func NewHandler(r *Repository) {
	Repo = r
}

// Index is the main page handler
func Index(w http.ResponseWriter, r *http.Request) {
	orders := Repo.Cache.LoadAll()

	data := make(map[string]interface{})
	data["orders"] = orders

	render.Template(
		w, r, "index.page.gohtml", &templatedata.TemplateData{
			Data: data,
		},
	)
}

// NotFound renders pages which not found
func NotFound(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "error.page.gohtml", &templatedata.TemplateData{})
}
