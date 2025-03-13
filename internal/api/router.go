package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rakheshkrishna2005/url-shortener/internal/api/handlers"
	"github.com/rakheshkrishna2005/url-shortener/internal/api/middleware"
)

// NewRouter sets up and configures the API router
func NewRouter(urlHandler *handlers.URLHandler, healthHandler *handlers.HealthHandler) *mux.Router {
	router := mux.NewRouter()

	// Apply common middleware
	router.Use(middleware.Logging)
	router.Use(middleware.Recovery)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
	
	// URL endpoints
	urlsRouter := api.PathPrefix("/urls").Subrouter()
	urlsRouter.HandleFunc("", urlHandler.CreateURL).Methods(http.MethodPost)
	urlsRouter.HandleFunc("/{id:[0-9]+}", urlHandler.GetURLByID).Methods(http.MethodGet)
	urlsRouter.HandleFunc("/{id:[0-9]+}", urlHandler.UpdateURL).Methods(http.MethodPut)
	urlsRouter.HandleFunc("/{id:[0-9]+}", urlHandler.DeleteURL).Methods(http.MethodDelete)

	// Health check
	router.HandleFunc("/health", healthHandler.HealthCheck).Methods(http.MethodGet)

	// Redirect handler
	router.HandleFunc("/{shortCode:[a-zA-Z0-9]+}", urlHandler.RedirectURL).Methods(http.MethodGet)

	// Serve static files and home page
	fs := http.FileServer(http.Dir("./web/static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	
	// Home page handler - we'll add this later
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./web/templates/index.html")
	}).Methods(http.MethodGet)

	return router
}