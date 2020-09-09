/* Quick start package for golang servers using Chi as router

Note: Must setup gotils.Logger first.

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
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/treeder/gotils"
	"go.uber.org/zap"
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
	return r
}

func Start(ctx context.Context, port int, r chi.Router) error {
	gotils.L(ctx).Sugar().Infof("Starting API server on port %v", port)
	srv := http.Server{Addr: fmt.Sprintf("0.0.0.0:%v", port), Handler: chi.ServerBaseContext(ctx, r)}
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		gotils.L(ctx).Error("error in http server", zap.Error(err))
	}
	return err
}

func SetupCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := gotils.WithRequestID(r.Context())
		ctx = gotils.AddFields(ctx, zap.String("path", r.URL.Path))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func Recoverer(next http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {

// 				logEntry := GetLogEntry(r)
// 				if logEntry != nil {
// 					logEntry.Panic(rvr, debug.Stack())
// 				} else {
// 					PrintPrettyStack(rvr)
// 				}

// 				w.WriteHeader(http.StatusInternalServerError)
// 			}
// 		}()

// 		next.ServeHTTP(w, r)
// 	}

// 	return http.HandlerFunc(fn)
// }

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
