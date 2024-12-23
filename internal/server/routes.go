package server

import (
	"embed"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/nrmnqdds/gomaluum/templates"
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

	r.Get("/favicon.ico", ServeFavicon)
	r.Get("/static/*", ServeStaticFiles)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		Chain(w, r, templates.Home())
	})

	r.Get("/reference", s.ScalarReference)

	r.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Use(s.PasetoAuthenticator())

			r.Post("/login", s.LoginHandler)
			r.Get("/schedule", s.ScheduleHandler)
			r.Get("/result", s.ResultHandler)
		})
	})

	return r
}

func ServeFavicon(w http.ResponseWriter, r *http.Request) {
	filePath := "favicon.ico"
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}

func ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path[len("/static/"):]
	fullPath := filepath.Join(".", "static", filePath)
	http.ServeFile(w, r, fullPath)
}
