/*

Quick start package for golang servers using Chi as router

Example:

	r := goapibase.InitRouter(ctx)
	// Setup your routes
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	// Start server
	_ = r.Start()

*/
package goapibase

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rs/cors"
	"github.com/treeder/gotils/v2"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func InitRouter(ctx context.Context) chi.Router {
	r := chi.NewRouter()
	r.Use(SetupCtx)
	r.Use(middleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           3600,
	}).Handler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		gotils.WriteError(w, http.StatusNotFound, gotils.ErrNotFound)
	})
	return r
}

// Regular server start
func Start(ctx context.Context, port int, r chi.Router) error {
	gotils.Logf(ctx, "info", "Starting API server on port %v", port)
	srv := http.Server{Addr: fmt.Sprintf("0.0.0.0:%v", port), Handler: r}
	srv.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		gotils.LogBeta(ctx, "error", "error in http.ListenAndServe: %v", err)
	}
	return err
}

// Starts with H2C support, should use when load balancer terminates the connection
// More info: https://cloud.google.com/run/docs/configuring/http2
func StartH2C(ctx context.Context, port int, r chi.Router) error {
	gotils.Logf(ctx, "info", "Starting API server on port %v", port)
	h2s := &http2.Server{} // http2 upgrades: https://www.mailgun.com/blog/http-2-cleartext-h2c-client-example-go/
	srv := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
		Handler: h2c.NewHandler(r, h2s), // HTTP/2 Cleartext handler
	}
	srv.BaseContext = func(_ net.Listener) context.Context {
		return ctx
	}
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		gotils.LogBeta(ctx, "error", "error in http.ListenAndServe: %v", err)
	}
	return err
}

func SetupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		id, _ := gonanoid.New()
		ctx = gotils.With(ctx, "requestID", id)
		ctx = gotils.With(ctx, "path", r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// WithValue is a middleware that sets a given key/value in a context chain.
func WithValue(key interface{}, val interface{}) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), key, val))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
