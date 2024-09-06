package routes

import (
	"github.com/go-chi/chi/v5"

	"git.sr.ht/~rehandaphedar/minv-server/handlers"
	"git.sr.ht/~rehandaphedar/minv-server/middlewares"
)

func setupChannelRoutes(r *chi.Mux) {
	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/channels", handlers.SelectChannels)
		r.Get("/channel/{channelname}", handlers.SelectChannel)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
		r.Get("/channel", handlers.SelectLoggedInChannel)
		r.Delete("/channel", handlers.DeleteChannel)
	})
}
