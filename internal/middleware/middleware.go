package middleware

import (
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/yalagtyarzh/L0/internal/config"
)

var app *config.AppConfig

// NewMiddleware sets the config for the middleware package
func NewMiddleware(a *config.AppConfig) {
	app = a
}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(
		http.Cookie{
			HttpOnly: true,
			Path:     "/",
			Secure:   app.InProduction,
			SameSite: http.SameSiteLaxMode,
		},
	)
	return csrfHandler
}
