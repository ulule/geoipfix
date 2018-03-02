package ipfix

import (
	"time"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"github.com/go-chi/chi/middleware"
)

type middlewareHandler = func(next http.Handler) http.Handler

type recoverMiddleware struct {
	debug bool
	logger *zap.Logger
}

func newRecoverMiddleware(debug bool, logger *zap.Logger) *recoverMiddleware {
	return &recoverMiddleware{
		debug: debug,
		logger: logger,
	}
}

func (m *recoverMiddleware) Handle(it interface{}) {
	if m.debug {
		fmt.Fprintf(os.Stderr, "Panic: %+v\n", it)
		debug.PrintStack()
	} else {
		m.logger.Error("Error handled")
	}
}

func (m *recoverMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				m.Handle(rvr)

				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func newLoggerMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

type StructuredLogger struct {
	Logger *zap.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: l.Logger}

	fields := []zapcore.Field{zap.String("ts", time.Now().UTC().Format(time.RFC1123))}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		fields = append(fields, zap.String("req.id", reqID))
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	fields = append(fields, []zapcore.Field{
		zap.String("http.scheme", scheme),
		zap.String("http.proto", r.Proto),
		zap.String("http.method", r.Method),
		zap.String("remote_addr", r.RemoteAddr),
		zap.String("user_agent", r.UserAgent()),
		zap.String("uri", fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)),
	}...)

	entry.Logger = l.Logger.With(fields...)

	entry.Logger.Info("request started")

	return entry
}

type StructuredLoggerEntry struct {
	Logger *zap.Logger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.With(
		zap.Int("res.status", status),
		zap.Int("res.bytes_length", bytes),
		zap.Float64("res.elapsed_ms", float64(elapsed.Nanoseconds()) / 1000000.0))

	l.Logger.Info("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.With(
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}
