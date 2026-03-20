package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sanjasimiccc/movieManagementSystem/middleware"
	"github.com/sanjasimiccc/movieManagementSystem/pkg/controllers"
)

type APIServer struct {
	addr    string
	handler *controllers.MovieHandler
}

func NewAPIServer(addr string, handler *controllers.MovieHandler) *APIServer {
	return &APIServer{
		addr:    addr,
		handler: handler,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	s.handler.RegisterRoutes(subrouter)

	return http.ListenAndServe(s.addr, middleware.APIKeyMiddleware(router))
}
