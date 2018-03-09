package geoipfix

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type handler func(w http.ResponseWriter, r *http.Request) error

// httpServer is an HTTP server.
type httpServer struct {
	srv     http.Server
	cfg     serverHTTPConfig
	mux     *chi.Mux
	opt     options
	recover *recoverMiddleware
}

// newHTTPServer retrieves a new HTTPServer instance.
func newHTTPServer(cfg serverHTTPConfig, opts ...option) *httpServer {
	opt := newOptions(opts...)

	srv := &httpServer{
		cfg: cfg,
	}

	opt.Logger = opt.Logger.With(zap.String("server", srv.Name()))

	srv.opt = opt

	return srv
}

func (h *httpServer) Name() string {
	return "http"
}

// handle handles an handler and captures error.
func (h *httpServer) handle(f handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			h.recover.Handle(err)
		}
	}
}

// Init initializes http server instance.
func (h *httpServer) Init() error {
	r := chi.NewRouter()

	h.recover = newRecoverMiddleware(h.opt.Debug, h.opt.Logger)

	cors := cors.New(cors.Options{
		AllowedOrigins:   h.cfg.Cors.AllowedOrigins,
		AllowedMethods:   h.cfg.Cors.AllowedMethods,
		AllowedHeaders:   h.cfg.Cors.AllowedHeaders,
		ExposedHeaders:   h.cfg.Cors.ExposedHeaders,
		AllowCredentials: h.cfg.Cors.AllowCredentials,
		MaxAge:           h.cfg.Cors.MaxAge,
	})
	r.Use(h.recover.Handler)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)
	r.Use(newLoggerMiddleware(h.opt.Logger))

	r.Get("/sys/health", func(w http.ResponseWriter, r *http.Request) {
		render.DefaultResponder(w, r, render.M{
			"status":     "OK",
			"version":    Version,
			"revision":   Revision,
			"build_time": BuildTime,
			"compiler":   Compiler,
		})
	})

	handler := httpHandler{h.opt}

	r.Get("/json/{ipAddress}", h.handle(handler.GetLocation))
	r.Get("/json/", h.handle(handler.GetLocation))

	h.mux = r

	return nil
}

// Serve serves http requests.
func (h *httpServer) Serve(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", strconv.Itoa(h.cfg.Port))

	h.srv = http.Server{
		Addr:    addr,
		Handler: chi.ServerBaseContext(ctx, h.mux),
	}
	h.opt.Logger.Info("Launch server", zap.String("addr", addr))

	return h.srv.ListenAndServe()
}

// Shutdown stops the http server.
func (h *httpServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := h.srv.Shutdown(ctx)

	h.opt.Logger.Info("Server shutdown")

	return err
}
