package app

import (
	"encoding/json"
	"net/http"

	"github.com/jonathanpopham/estate-agent/internal/config"
)

type Server struct {
	cfg config.Config
	mux *http.ServeMux
}

func NewServer(cfg config.Config) http.Handler {
	s := &Server{
		cfg: cfg,
		mux: http.NewServeMux(),
	}
	s.routes()
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"cloud":   s.cfg.Cloud,
		"dry_run": s.cfg.DryRun,
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
