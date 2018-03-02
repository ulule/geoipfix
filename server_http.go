package ipfix

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/valve"
	"go.uber.org/zap"
)

// HTTPServer is an HTTP server.
type HTTPServer struct {
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
		w.Write([]byte("OK"))
	})
	r.Get("/json/{ipAddress}", h.handle(IPAddressHandler))
	r.Get("/json/", h.handle(IPAddressHandler))

	h.mux = r

	return nil
}

// Serve serves http requests.
func (h *HTTPServer) Serve() error {
	addr := fmt.Sprintf(":%s", strconv.Itoa(h.cfg.Port))

	valv := valve.New()
	baseCtx := valv.Context()

	srv := http.Server{
		Addr:    addr,
		Handler: chi.ServerBaseContext(baseCtx, h.mux),
	}
	h.opt.Logger.Info("Launch HTTP server", zap.String("addr", addr))
	defer h.opt.Logger.Sync()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			h.opt.Logger.Info("shutting down..")

			// first valv
			valv.Shutdown(20 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// start http shutdown
			srv.Shutdown(ctx)

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				h.opt.Logger.Info("not all connections have been closed")
			case <-ctx.Done():

			}
		}
	}()
	return srv.ListenAndServe()
}
