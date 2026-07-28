package handlers

import "net/http"

func (s *Server) Home(w http.ResponseWriter, r *http.Request) {
	s.render.Render(w, http.StatusOK, "home.html", nil)
}

func (s *Server) NotFound(w http.ResponseWriter, r *http.Request) {
	s.render.Render(w, http.StatusNotFound, "404.html", nil)
}
