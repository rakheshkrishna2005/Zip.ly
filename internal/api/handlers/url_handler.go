package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rakheshkrishna2005/url-shortener/internal/models"
	"github.com/rakheshkrishna2005/url-shortener/internal/service"
)

// URLHandler handles HTTP requests for URL operations
type URLHandler struct {
	urlService *service.URLService
}

// NewURLHandler creates a new URLHandler
func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

// CreateURL handles POST requests to create a new short URL
func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {
	var req models.CreateURLRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Get user IP
	userIP := getUserIP(r)

	// Create short URL
	resp, err := h.urlService.CreateShortURL(r.Context(), req, userIP)
	if err != nil {
		if err == models.ErrDuplicateAlias {
			http.Error(w, "Custom alias already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to create short URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

// GetURLByID handles GET requests to retrieve URL details by ID
func (h *URLHandler) GetURLByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	details, err := h.urlService.GetURLByID(r.Context(), id)
	if err != nil {
		if err == models.ErrURLNotFound {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to retrieve URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(details)
}

// UpdateURL handles PUT requests to update URL details
func (h *URLHandler) UpdateURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	var req struct {
		OriginalURL string  `json:"original_url,omitempty"`
		CustomAlias *string `json:"custom_alias,omitempty"`
		ExpiresIn   *int    `json:"expires_in,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	err = h.urlService.UpdateURL(r.Context(), id, req.OriginalURL, req.CustomAlias, req.ExpiresIn)
	if err != nil {
		if err == models.ErrURLNotFound {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		} else if err == models.ErrDuplicateAlias {
			http.Error(w, "Custom alias already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Failed to update URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteURL handles DELETE requests to remove a URL
func (h *URLHandler) DeleteURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid URL ID", http.StatusBadRequest)
		return
	}

	err = h.urlService.DeleteURL(r.Context(), id)
	if err != nil {
		if err == models.ErrURLNotFound {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RedirectURL handles GET requests to redirect short URLs to their original destination
func (h *URLHandler) RedirectURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortCode := vars["shortCode"]

	url, err := h.urlService.GetURL(r.Context(), shortCode)
	if err != nil {
		if err == models.ErrURLNotFound {
			http.Error(w, "URL not found", http.StatusNotFound)
			return
		} else if err == models.ErrURLExpired {
			http.Error(w, "URL has expired", http.StatusGone)
			return
		} else if err == models.ErrInvalidShortCode {
			http.Error(w, "Invalid short code", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to retrieve URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Record click asynchronously to not delay the redirect
	go func() {
		referer := r.Header.Get("Referer")
		userAgent := r.Header.Get("User-Agent")
		ipAddress := getUserIP(r)
		// Using background context since request context will be canceled
		h.urlService.RecordClick(r.Context(), url.ID, referer, userAgent, ipAddress)
	}()

	http.Redirect(w, r, url.OriginalURL, http.StatusFound)
}

// getUserIP extracts the client IP address from request
func getUserIP(r *http.Request) string {
	// Check for X-Forwarded-For header first (for proxied requests)
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fall back to RemoteAddr
	ip = r.RemoteAddr
	// Strip port if present
	if i := strings.LastIndex(ip, ":"); i != -1 {
		ip = ip[:i]
	}
	return ip
}