package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/natecw/minily/models"
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
	router.HandleFunc("POST /{$}", s.create)
	router.HandleFunc("GET /{short_code}", s.redirect)
	return router
}

func (s *Server) create(w http.ResponseWriter, r *http.Request) {
	s.log.Info("request received to create short_code")
	ct := r.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediaType != "application/json" {
		http.Error(w, "expected json content-type", http.StatusUnsupportedMediaType)
		return
	}

	var request models.CreateRequest
	parser := json.NewDecoder(r.Body)
	parser.DisallowUnknownFields()
	if err := parser.Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.dumpRequest(r)
		return
	}

	if request.URL == "" {
		http.Error(w, "url is missing", http.StatusBadRequest)
		return
	}

	mini, err := s.storage.CreateMinily(r.Context(), request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.dumpRequest(r)
		return
	}
	s.renderJson(w, mini)
}

func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {
	s.log.Info("redirect request received", "path", r.URL.Path)
	short_code := r.PathValue("short_code")
	if short_code == "" {
		http.Error(w, "missing path", http.StatusBadRequest)
		return
	}

	long_url, err := s.storage.GetUrl(r.Context(), short_code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.log.Error("error retrieving", "short_code", short_code)
		return
	}
	s.log.Info("redirecting", "short_code", short_code, "location", long_url)
	http.Redirect(w, r, long_url, http.StatusTemporaryRedirect)
}

func (s *Server) dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		return
	}
	s.log.Error("error parsing json", "body", string(dump))
}

func (s *Server) renderJson(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		s.log.Error("rendering json", "error", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}
