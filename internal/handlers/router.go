package handlers

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/RobertLMcCrary/portfolio/internal/blog"
	"github.com/RobertLMcCrary/portfolio/internal/render"
)

type Server struct {
	render *render.Renderer
	posts  []blog.Post
}

func New(r *render.Renderer, posts []blog.Post) *Server {
	return &Server{render: r, posts: posts}
}

func (s *Server) Routes(static fs.FS) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", s.Home)
	mux.HandleFunc("GET /blog", s.BlogList)
	mux.HandleFunc("GET /blog/{slug}", s.BlogPost)
	mux.Handle("GET /static/", http.FileServerFS(static))
	mux.HandleFunc("/", s.NotFound)

	return withLogging(mux)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}
