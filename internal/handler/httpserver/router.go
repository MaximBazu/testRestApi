package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(userHandler *UserHandler, productHandler ProductHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(5 * time.Second))

	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.List)
		r.Post("/", userHandler.Create)
		r.Get("/{id}", userHandler.GetByID)
		r.Delete("/{id}", userHandler.Delete)
	})

	r.Route("/products", func(r chi.Router) {
		r.Post("/", productHandler.Create)
		r.Get("/", productHandler.List)
		r.Get("/{id}", productHandler.GetByID)
		r.Patch("/{id}", productHandler.Update)
		r.Delete("/{id}", productHandler.Delete)
	})

	return r
}
