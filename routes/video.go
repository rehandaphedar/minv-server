package routes

import (
	"github.com/go-chi/chi/v5"

	"git.sr.ht/~rehandaphedar/minv-server/handlers"
	"git.sr.ht/~rehandaphedar/minv-server/middlewares"
)

func setupVideoRoutes(r *chi.Mux) {
	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/videos", handlers.SelectVideos)
		r.Get("/video/{slug}", handlers.SelectVideo)
		r.Get("/videos/byChannel/{channelname}", handlers.SelectVideosByChannel)
	})

	// Private routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware)

		r.Post("/videos", handlers.CreateVideo)
		r.Delete("/video/{slug}", handlers.DeleteVideo)
		r.Put("/video/{slug}", handlers.UpdateVideo)
	})
}
