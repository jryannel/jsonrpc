package jsonrpc

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

type Server struct {
	r chi.Router
}

func NewHTTPServer() *Server {
	r := NewRouter()
	return &Server{
		r: r,
	}
}

func (s *Server) Router() chi.Router {
	return s.r
}

func (s *Server) Start(addr string) error {
	log.Printf("start http server at %s", addr)
	return http.ListenAndServe(addr, s.r)
}
