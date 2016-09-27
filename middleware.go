package main

import (
	"fmt"
	"net/http"
	"time"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
)

// Middleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type LogrusMiddleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *logrus.Logger
	// Name is the name of the application as recorded in latency metrics
	Name string
}

// NewMiddleware returns a new *Middleware, yay!
func NewLogrusMiddleware() *LogrusMiddleware {
	return NewCustomLogrusMiddleware(logrus.InfoLevel, &logrus.TextFormatter{}, "web")
}

// NewCustomMiddleware builds a *Middleware with the given level and formatter
func NewCustomLogrusMiddleware(level logrus.Level, formatter logrus.Formatter, name string) *LogrusMiddleware {
	log := logrus.New()
	log.Level = level
	log.Formatter = formatter

	return &LogrusMiddleware{Logger: log, Name: name}
}

func (l *LogrusMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	headers := []string{}

	for k, v := range r.Header {
		headers = append(headers, fmt.Sprintf("%s: %s", k, v))
	}

	l.Logger.WithFields(logrus.Fields{
		"method":  r.Method,
		"request": r.RequestURI,
		"remote":  r.RemoteAddr,
		"headers": strings.Join(headers, " | "),
	}).Info("started handling request")

	next(rw, r)

	latency := time.Since(start)
	res := rw.(negroni.ResponseWriter)
	l.Logger.WithFields(logrus.Fields{
		"status":      res.Status(),
		"method":      r.Method,
		"request":     r.RequestURI,
		"remote":      r.RemoteAddr,
		"text_status": http.StatusText(res.Status()),
		"took":        latency,
		fmt.Sprintf("measure#%s.latency", l.Name): latency.Nanoseconds(),
	}).Info("completed handling request")
}
