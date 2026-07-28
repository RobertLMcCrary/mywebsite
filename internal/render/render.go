package render

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

// Pages that pair with templates/layout.html. Each must define a
// "content" template block for the layout to embed.
var pageNames = []string{"home.html", "blog_list.html", "blog_post.html", "404.html"}

var funcs = template.FuncMap{
	"now": time.Now,
}

type Renderer struct {
	pages map[string]*template.Template
}

func New(fsys embed.FS) (*Renderer, error) {
	pages := make(map[string]*template.Template, len(pageNames))
	for _, name := range pageNames {
		tmpl, err := template.New("layout.html").Funcs(funcs).ParseFS(fsys, "templates/layout.html", "templates/"+name)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", name, err)
		}
		pages[name] = tmpl
	}
	return &Renderer{pages: pages}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, status int, name string, data any) {
	tmpl, ok := r.pages[name]
	if !ok {
		http.Error(w, fmt.Sprintf("unknown template %q", name), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := tmpl.ExecuteTemplate(w, "layout.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
