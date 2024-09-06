package routes

import (
	"git.sr.ht/~rehandaphedar/minv-server/handlers"
	"git.sr.ht/~rehandaphedar/minv-server/middlewares"
	"github.com/go-chi/chi/v5"
)

func setupAtomRoutes(r *chi.Mux) {
	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/videos.atom", handlers.AtomSelectVideos)
		r.Get("/channels.atom", handlers.AtomSelectChannels)
		r.Get("/channel/{channelname}.atom", handlers.AtomSelectVideosByChannel)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)
	})
}
