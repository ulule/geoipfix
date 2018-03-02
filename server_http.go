package ipfix

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

// HTTPServer is an HTTP server.
type HTTPServer struct {
	srv     http.Server
	cfg     serverHTTPConfig
	mux     *chi.Mux
	opt     Options
	recover *Recover
}

// NewHTTPServer retrieves a new HTTPServer instance.
func NewHTTPServer(cfg serverHTTPConfig, opts ...Option) *HTTPServer {
	opt := NewOptions(opts...)

	return &HTTPServer{
		cfg: cfg,
		opt: opt,
	}
}

func (h *HTTPServer) handle(f Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(h.opt, w, r)
		if err != nil {
			h.recover.Handle(err)
		}
	}
}

// Init initializes HTTPServer instance.
func (h *HTTPServer) Init() error {
	r := chi.NewRouter()

	h.recover = NewRecover(h.opt.Debug, h.opt.Logger)

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

	r.Get("/sys/health", func(w http.ResponseWriter, r *http.Request) {
		render.DefaultResponder(w, r, render.M{
			"status":     "OK",
			"version":    Version,
			"revision":   Revision,
			"build_time": BuildTime,
			"compiler":   Compiler,
		})
	})
	r.Get("/json/{ipAddress}", h.handle(IPAddressHandler))
	r.Get("/json/", h.handle(IPAddressHandler))

	h.mux = r

	return nil
}

// Serve serves http requests.
func (h *HTTPServer) Serve(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", strconv.Itoa(h.cfg.Port))

	h.srv = http.Server{
		Addr:    addr,
		Handler: chi.ServerBaseContext(ctx, h.mux),
	}
	h.opt.Logger.Info("Launch HTTP server", zap.String("addr", addr))

	return h.srv.ListenAndServe()
}

func (h *HTTPServer) Shutdown() {
	// create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// start http shutdown
	h.srv.Shutdown(ctx)

	h.opt.Logger.Info("HTTP server successfully shutdown")
}
