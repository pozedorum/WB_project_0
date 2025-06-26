package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// API endpoint
	mux.HandleFunc("/api/order/", s.orderHandler)

	// Главная страница
	mux.HandleFunc("/", s.indexHandler)

	return mux
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := s.frontend.RenderIndex(w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) orderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid := r.URL.Path[len("/order/"):]
	if uid == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := s.cache.Get(uid)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
