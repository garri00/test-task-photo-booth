package middleware

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/spf13/viper"
)

const (
	maxAge = 300
)

// NewDefaultCors set default cors params
func NewDefaultCors(r *chi.Mux) {
	//TODO setup cors properly
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// AllowedOrigins: []string{"https://*", "http://*"}, //localhost for capacitor app
		AllowedOrigins: []string{"*"}, // Allow all origins
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		//AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           maxAge, // Maximum value not ignored by any of major browsers
		Debug:            viper.GetBool("middlewares.cors"),
	}))
}
