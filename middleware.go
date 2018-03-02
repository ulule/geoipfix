package ipfix

import (
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"go.uber.org/zap"
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
