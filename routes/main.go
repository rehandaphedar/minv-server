package routes

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

func Setup(r *chi.Mux) {

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: viper.GetStringSlice("allowed_origins"),
		AllowedMethods: []string{"HEAD", "GET", "POST", "DELETE", "PUT"},
		AllowCredentials: true,
	}))

	setupAuthRoutes(r)
	setupChannelRoutes(r)
	setupVideoRoutes(r)
	setupAtomRoutes(r)
}
