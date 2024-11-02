package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/natecw/minily/storage"
)

var defaultStopTimeout = time.Second * 30

type Server struct {
	addr    string
	storage *storage.Storage
	log     *slog.Logger
}

func NewApi(addr string, storage *storage.Storage, logger *slog.Logger) (*Server, error) {
	if addr == "" {
		return nil, errors.New("addr cannot be blank")
	}

	return &Server{
		addr:    addr,
		storage: storage,
		log:     logger,
	}, nil
}

func (s *Server) Start(stop <-chan struct{}) error {
	srv := &http.Server{
		Addr:    s.addr,
		Handler: s.router(),
	}

	go func() {
		s.log.Info("starting server", "location", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("listen", "error", err)
			os.Exit(1)
		}
	}()

	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), defaultStopTimeout)
	defer cancel()

	s.log.Info("stopping server", "timeout", defaultStopTimeout)
	return srv.Shutdown(ctx)
}

func (s *Server) router() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("POST /{$}", s.Create)
	router.HandleFunc("GET /{short_code}", s.Redirect)
	return router
}
