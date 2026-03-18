package main

import (
	"curso-gitops/internal/auth"
	"curso-gitops/internal/handlers"
	"curso-gitops/internal/repository"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := repository.ConnectDB(); err != nil {
		log.Fatalf("Error conectando a la base de datos: %v", err)
	}

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(auth.SecurityHeaders)
	r.Use(middleware.Throttle(100)) // Max 100 concurrent requests

	r.Route("/api", func(r chi.Router) {
		r.Post("/login", handlers.Login)
		r.Post("/register", handlers.Register)

		r.Group(func(r chi.Router) {
			r.Use(auth.JWTMiddleware)
			r.Get("/modules", handlers.GetModules)

			r.Group(func(r chi.Router) {
				r.Use(auth.AdminMiddleware)
				r.Post("/modules", handlers.CreateModule)
				r.Delete("/modules/{id}", handlers.DeleteModule)
				r.Put("/modules/{id}/toggle", handlers.ToggleModule)
				r.Put("/modules/{id}", handlers.UpdateModule)

				r.Post("/episodes", handlers.CreateEpisode)
				r.Delete("/episodes/{id}", handlers.DeleteEpisode)
				r.Put("/episodes/{id}", handlers.UpdateEpisode)
			})
		})
	})

	workDir, _ := os.Getwd()
	frontendDir := filepath.Join(workDir, "frontend")

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontendDir, "index.html"))
	})
	r.Get("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontendDir, "dashboard.html"))
	})
	r.Get("/admin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(frontendDir, "admin.html"))
	})

	r.Handle("/*", http.FileServer(http.Dir(frontendDir)))

	log.Println("Servidor iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
