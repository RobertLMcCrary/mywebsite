package blog

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"gopkg.in/yaml.v3"
)

type frontmatter struct {
	Title   string `yaml:"title"`
	Date    string `yaml:"date"`
	Slug    string `yaml:"slug"`
	Summary string `yaml:"summary"`
}

// Load reads every *.md file in dir (within fsys), parses its YAML
// frontmatter and markdown body, and returns the posts sorted by date
// descending (newest first).
func Load(fsys fs.FS, dir string) ([]Post, error) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return nil, fmt.Errorf("reading posts dir: %w", err)
	}

	var posts []Post
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		raw, err := fs.ReadFile(fsys, filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", entry.Name(), err)
		}

		post, err := parsePost(entry.Name(), raw)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", entry.Name(), err)
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts, nil
}

func parsePost(filename string, raw []byte) (Post, error) {
	fm, body, err := splitFrontmatter(raw)
	if err != nil {
		return Post{}, err
	}

	var meta frontmatter
	if err := yaml.Unmarshal(fm, &meta); err != nil {
		return Post{}, fmt.Errorf("parsing frontmatter: %w", err)
	}

	derivedDate, derivedSlug := splitFilename(filename)
	slug := meta.Slug
	if slug == "" {
		slug = derivedSlug
	}
	date := meta.Date
	if date == "" {
		date = derivedDate
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return Post{}, fmt.Errorf("parsing date %q: %w", date, err)
	}

	var buf bytes.Buffer
	if err := goldmark.Convert(body, &buf); err != nil {
		return Post{}, fmt.Errorf("rendering markdown: %w", err)
	}

	return Post{
		Title:   meta.Title,
		Slug:    slug,
		Date:    parsedDate,
		Summary: meta.Summary,
		Body:    template.HTML(buf.String()),
	}, nil
}

// splitFrontmatter separates a leading "---" delimited YAML block from
// the remaining markdown body. If no frontmatter block is present, fm
// is nil and the entire input is treated as body.
func splitFrontmatter(raw []byte) (fm []byte, body []byte, err error) {
	const delim = "---"

	trimmed := bytes.TrimLeft(raw, "\n")
	if !bytes.HasPrefix(trimmed, []byte(delim)) {
		return nil, raw, nil
	}

	rest := trimmed[len(delim):]
	fm, body, ok := bytes.Cut(rest, []byte("\n"+delim))
	if !ok {
		return nil, nil, fmt.Errorf("unterminated frontmatter")
	}

	return fm, body, nil
}

// splitFilename derives a date and slug from a "YYYY-MM-DD-slug.md"
// filename, used as a fallback when frontmatter omits either field.
func splitFilename(filename string) (date, slug string) {
	name := strings.TrimSuffix(filename, filepath.Ext(filename))
	parts := strings.SplitN(name, "-", 4)
	if len(parts) == 4 {
		return strings.Join(parts[:3], "-"), parts[3]
	}
	return "", name
}
