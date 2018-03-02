package ipfix

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fiorix/freegeoip"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

// HTTPServer is an HTTP server.
type HTTPServer struct {
	cfg serverHTTPConfig
	mux *chi.Mux
	opt Options
}

// NewHTTPServer retrieves a new HTTPServer instance.
func NewHTTPServer(cfg serverHTTPConfig, opts ...Option) *HTTPServer {
	opt := NewOptions(opts...)

	return &HTTPServer{
		cfg: cfg,
		opt: opt,
	}
}

// Init initializes HTTPServer instance.
func (h *HTTPServer) Init() error {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedOrigins:   h.cfg.Cors.AllowedOrigins,
		AllowedMethods:   h.cfg.Cors.AllowedMethods,
		AllowedHeaders:   h.cfg.Cors.AllowedHeaders,
		ExposedHeaders:   h.cfg.Cors.ExposedHeaders,
		AllowCredentials: h.cfg.Cors.AllowCredentials,
		MaxAge:           h.cfg.Cors.MaxAge,
	})
	r.Use(cors.Handler)
	r.Use(middleware.RequestID)

	geoipHandler := freegeoip.ProxyHandler(freegeoip.NewHandler(h.opt.DB, &freegeoip.JSONEncoder{}))
	handler := func(w http.ResponseWriter, r *http.Request) {
		geoipHandler.ServeHTTP(w, r)
	}
	r.Get("/sys/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	r.Get("/json/{ipAddress}", handler)
	r.Get("/json/", handler)

	h.mux = r

	return nil
}

// Serve serves http requests.
func (h *HTTPServer) Serve() error {
	addr := fmt.Sprintf(":%s", strconv.Itoa(h.cfg.Port))

	h.opt.Logger.Info("Launch HTTP server", zap.String("addr", addr))
	defer h.opt.Logger.Sync()

	return http.ListenAndServe(addr, h.mux)
}
