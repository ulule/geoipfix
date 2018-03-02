package ipfix

import (
	"github.com/fiorix/freegeoip"
	"go.uber.org/zap"
)

// option is a functional option.
type option func(*options)

// NewOptions initializes ipfix options.
func newOptions(opts ...option) options {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// options are ipfix options.
type options struct {
	DB     *freegeoip.DB
	Logger *zap.Logger
	Debug  bool
}

// withDB sets the database.
func withDB(db *freegeoip.DB) option {
	return func(o *options) {
		o.DB = db
	}
}

// withDebug sets the debug flag.
func withDebug(debug bool) option {
	return func(o *options) {
		o.Debug = debug
	}
}

// withLogger sets the logger.
func withLogger(logger *zap.Logger) option {
	return func(o *options) {
		o.Logger = logger
	}
}
