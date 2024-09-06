package routes

import (
	"github.com/go-chi/chi/v5"

	"git.sr.ht/~rehandaphedar/minv-server/handlers"
	"git.sr.ht/~rehandaphedar/minv-server/middlewares"
)

func setupAuthRoutes(r *chi.Mux) {
	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/auth", handlers.Auth)
		r.Post("/logout", handlers.Logout)
	})

	// Private routes
	r.Group(func(r chi.Router) {

		r.Use(middlewares.AuthMiddleware)
	})
}
