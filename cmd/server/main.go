package main
import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	assets "github.com/RobertLMcCrary/portfolio"
	"github.com/RobertLMcCrary/portfolio/internal/blog"
	"github.com/RobertLMcCrary/portfolio/internal/handlers"
	"github.com/RobertLMcCrary/portfolio/internal/render"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	postsDir := envOr("POSTS_DIR", "./content/posts")
	posts, err := blog.Load(os.DirFS(postsDir), ".")
	if err != nil {
		return err
	}

	renderer, err := render.New(assets.Templates)
	if err != nil {
		return err
	}

	server := handlers.New(renderer, posts)

	httpServer := &http.Server{
		Addr:         envOr("LISTEN_ADDR", ":3000"),
		Handler:      server.Routes(assets.Static),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return httpServer.Shutdown(shutdownCtx)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
