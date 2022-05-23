package mux

import (
	"github.com/go-chi/chi"
)

func Router() *chi.Mux {
	mux := chi.NewRouter()

	return mux
}
