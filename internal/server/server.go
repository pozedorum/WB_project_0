// Package server is responcible for how the server works
package server

// import (
// 	"context"
// 	"log"
// 	"net/http"
// 	"time"

// 	"wb_project_0/internal/database"
// 	"wb_project_0/internal/middleware"
// )

// type Server struct {
// 	httpServer *http.Server
// 	db         *database.Database
// }

// func New(dbs *database.Database, addr string) *Server {
// 	mux := http.NewServeMux()

// 	return &Server{
// 		httpServer: &http.Server{
// 			Addr:         addr,
// 			Handler:      mux,
// 			ReadTimeout:  10 * time.Second,
// 			WriteTimeout: 10 * time.Second,
// 		},
// 		db: dbs,
// 	}
// }

// func (s *Server) RegisterHandlers() {
// 	mux := s.httpServer.Handler.(*http.ServeMux)

// 	// Публичные маршруты
// 	mux.HandleFunc("/", HomeHandler(s.db))
// 	mux.HandleFunc("/article/", ArticleHandler(s.db))

// 	// Закрытые маршруты (с middleware)
// 	adminHandler := http.HandlerFunc(AdminHandler(s.db))
// 	mux.Handle("/admin", middleware.AdminAuth(adminHandler))
// }

// func (s *Server) Start() error {
// 	log.Printf("Server started on %s", s.httpServer.Addr)
// 	return s.httpServer.ListenAndServe()
// }

// func (s *Server) Shutdown(ctx context.Context) error {
// 	return s.httpServer.Shutdown(ctx)
// }
