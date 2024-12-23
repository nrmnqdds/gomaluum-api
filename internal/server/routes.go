package server

import (
	"embed"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var DocsPath embed.FS

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		// MaxAge:           300,
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.RedirectSlashes)

	// redirect to /reference
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/reference", http.StatusMovedPermanently)
	})

	r.Get("/reference", s.ScalarReference)

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", s.LoginHandler)

		r.Group(func(r chi.Router) {
			r.Use(s.PasetoAuthenticator())

			r.Get("/schedule", s.ScheduleHandler)
			r.Get("/result", s.ResultHandler)
		})
	})

	return r
}
