package server

import (
	"encoding/json"
	"net/http"
	"strings"
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
		s.logger.Printf("error: Wrong method on orderHandler")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	uid, found := strings.CutPrefix(r.URL.Path, "/api/order/")
	if uid == "" || found == false {
		s.logger.Printf("error: No uid in request")
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}
	s.logger.Printf("Recieved uid from request: %s", uid)
	order, err := s.cache.Get(uid)
	if err != nil {
		s.logger.Printf("error: Order not found")
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}
	s.logger.Printf("order found: %s", order.OrderUID)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(order); err != nil {
		s.logger.Printf("error: Failed to encode responce")
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
