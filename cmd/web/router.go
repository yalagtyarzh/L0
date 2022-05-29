package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/yalagtyarzh/L0/internal/handlers"
)

// Router returns new http handler with routers
func Router() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)

	mux.NotFound(handlers.NotFound)
	mux.Get("/", handlers.Index)
	mux.Get("/{id}", handlers.ShowOrder)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
