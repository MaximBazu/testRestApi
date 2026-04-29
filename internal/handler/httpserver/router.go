package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	userHandler *UserHandler,
	productHandler *ProductHandler,
	orderHandler *OrderHandler,
	orderItemHandler *OrderItemHandler,
	productSizeHandler *ProductSizeHandler,
	productImageHandler *ProductImageHandler,
) http.Handler {
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

	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.Create)
		r.Get("/", orderHandler.List)
		r.Get("/{id}", orderHandler.GetByID)
		r.Delete("/{id}", orderHandler.Delete)
	})

	r.Route("/order-items", func(r chi.Router) {
		r.Post("/", orderItemHandler.Create)
		r.Get("/", orderItemHandler.List)
		r.Get("/{id}", orderItemHandler.GetByID)
		r.Delete("/{id}", orderItemHandler.Delete)
	})

	r.Route("/product-sizes", func(r chi.Router) {
		r.Post("/", productSizeHandler.Create)
		r.Get("/", productSizeHandler.List)
		r.Get("/{id}", productSizeHandler.GetByID)
		r.Delete("/{id}", productSizeHandler.Delete)
	})

	r.Route("/product-images", func(r chi.Router) {
		r.Post("/", productImageHandler.Create)
		r.Get("/", productImageHandler.List)
		r.Get("/one", productImageHandler.GetByKey)
		r.Delete("/", productImageHandler.Delete)
	})

	return r
}
