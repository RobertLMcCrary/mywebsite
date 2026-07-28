package handlers

import (
	"net/http"

	"github.com/RobertLMcCrary/portfolio/internal/blog"
)

type blogListData struct {
	Posts []blog.Post
}

func (s *Server) BlogList(w http.ResponseWriter, r *http.Request) {
	s.render.Render(w, http.StatusOK, "blog_list.html", blogListData{Posts: s.posts})
}

func (s *Server) BlogPost(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	post, ok := blog.FindBySlug(s.posts, slug)
	if !ok {
		s.NotFound(w, r)
		return
	}

	s.render.Render(w, http.StatusOK, "blog_post.html", post)
}
