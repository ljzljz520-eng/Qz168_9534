package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"tea17/internal/catalog"
	"tea17/internal/domain"
	"tea17/internal/service"
	"tea17/internal/store"
)

type Server struct {
	service *service.Service
	mux     *http.ServeMux
}

func New(s *service.Service) *Server {
	server := &Server{service: s, mux: http.NewServeMux()}
	server.routes()
	return server
}
func (s *Server) Handler() http.Handler { return s.mux }
func (s *Server) routes() {
	s.mux.HandleFunc("GET /health", s.health)
	s.mux.HandleFunc("POST /records", s.createRecord)
	s.mux.HandleFunc("GET /records/{id}", s.getRecord)
	s.mux.HandleFunc("POST /records/{id}/reviews", s.reviewRecord)
	s.mux.HandleFunc("GET /records/{id}/reviews", s.reviewHistory)
	s.mux.HandleFunc("POST /profiles", s.createProfile)
	s.mux.HandleFunc("GET /profiles/{id}/recommendations", s.recommendations)
	s.mux.HandleFunc("GET /templates", s.templates)
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "tea17"})
}
func (s *Server) createRecord(w http.ResponseWriter, r *http.Request) {
	var input domain.BeverageRecord
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	record, err := s.service.Register(input, r.Header.Get("X-Actor-ID"))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, record)
}
func (s *Server) getRecord(w http.ResponseWriter, r *http.Request) {
	record, err := s.service.Record(r.PathValue("id"))
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, store.ErrNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, record)
}
func (s *Server) reviewRecord(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ViewerID string `json:"viewer_id"`
		Decision string `json:"decision"`
		Note     string `json:"note"`
	}
	if err := decode(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	review, err := s.service.Review(service.ReviewInput{RecordID: r.PathValue("id"), ViewerID: input.ViewerID, Decision: input.Decision, Note: input.Note})
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, review)
}
func (s *Server) reviewHistory(w http.ResponseWriter, r *http.Request) {
	items, err := s.service.History(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (s *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	var profile domain.CustomerProfile
	if err := decode(r, &profile); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := s.service.SubmitProfile(profile); err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusCreated, profile)
}
func (s *Server) recommendations(w http.ResponseWriter, r *http.Request) {
	limit := 5
	if raw := r.URL.Query().Get("limit"); raw != "" {
		value, err := strconv.Atoi(raw)
		if err != nil || value < 1 || value > 50 {
			writeError(w, http.StatusBadRequest, errors.New("limit must be between 1 and 50"))
			return
		}
		limit = value
	}
	items, err := s.service.Recommendations(r.PathValue("id"), limit)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err)
		return
	}
	writeJSON(w, http.StatusOK, items)
}
func (s *Server) templates(w http.ResponseWriter, r *http.Request) {
	maxPrice := 0
	if raw := r.URL.Query().Get("max_price"); raw != "" {
		maxPrice, _ = strconv.Atoi(raw)
	}
	writeJSON(w, http.StatusOK, catalog.SearchTemplates(strings.TrimSpace(r.URL.Query().Get("q")), maxPrice))
}
