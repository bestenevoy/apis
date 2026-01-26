package server

import (
	"net/http"

	"wrzapi/internal/handlers"
)

type Server struct {
	mux *http.ServeMux
}

func New() *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handlers.Health)
	mux.HandleFunc("/api/time", handlers.Time)
	mux.HandleFunc("/api/page-info", handlers.PageInfo)
	mux.HandleFunc("/openapi.yaml", handlers.OpenAPI)
	mux.HandleFunc("/docs", handlers.Docs)

	return &Server{mux: mux}
}

func (s *Server) ListenAndServe(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	return srv.ListenAndServe()
}