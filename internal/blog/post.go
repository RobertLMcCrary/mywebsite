package blog

import (
	"html/template"
	"time"
)

type Post struct {
	Title   string
	Slug    string
	Date    time.Time
	Summary string
	Body    template.HTML
}

func FindBySlug(posts []Post, slug string) (Post, bool) {
	for _, p := range posts {
		if p.Slug == slug {
			return p, true
		}
	}
	return Post{}, false
}
