package api

import (
	"encoding/json"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/natecw/minily/models"
)

func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
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

	_, err = url.Parse(request.URL)
	if err != nil {
		http.Error(w, "invalid 'long_url' given", http.StatusBadRequest)
		s.log.Error("invalid url", "url", request.URL, "err", err.Error())
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

func (s *Server) Redirect(w http.ResponseWriter, r *http.Request) {
	s.log.Info("redirect request received", "path", r.URL.Path)
	short_code := r.PathValue("short_code")
	if short_code == "" {
		http.Error(w, "missing path", http.StatusBadRequest)
		return
	}

	long_url, err := s.storage.GetOriginalUrl(r.Context(), short_code)
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
