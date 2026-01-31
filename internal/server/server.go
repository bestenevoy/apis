package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"wrzapi/internal/handlers"
	"wrzapi/nav"
)

type Server struct {
	engine *gin.Engine
}

type Config struct {
	NavDataPath string
	NavDev      bool
}

func New(cfg Config) (*Server, error) {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/healthz", handlers.Health)
	engine.GET("/api/time", handlers.Time)
	engine.GET("/api/page-info", handlers.PageInfo)
	engine.GET("/openapi.yaml", handlers.OpenAPI)
	engine.GET("/openapi.json", handlers.OpenAPIJSON)
	engine.GET("/docs", handlers.Docs)

	navApp, err := nav.New(nav.Config{
		DataPath: cfg.NavDataPath,
		Dev:      cfg.NavDev,
	})
	if err != nil {
		return nil, err
	}
	engine.NoRoute(gin.WrapH(navApp.Handler()))

	return &Server{engine: engine}, nil
}

func (s *Server) ListenAndServe(addr string) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: s.engine,
	}
	return srv.ListenAndServe()
}
