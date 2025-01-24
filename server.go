package recordsrestapi

import (
	"context"
	"net/http"
	"time"
)

type Server struct {
	http *http.Server
}

func (s *Server) Start(port string, handler http.Handler) error {
	s.http = &http.Server{
		Addr:    ":" + port,
		Handler: handler,
		ReadTimeout: 10* time.Second,
		WriteTimeout: 10* time.Second,
	}
	return s.http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}