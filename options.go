package ipfix

import (
	"github.com/fiorix/freegeoip"
	"go.uber.org/zap"
)

// Option is a functional option.
type Option func(*Options)

// NewOptions initializes ipfix options.
func NewOptions(opts ...Option) Options {
	opt := Options{}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}

// Options are ipfix options.
type Options struct {
	DB     *freegeoip.DB
	Logger *zap.Logger
	Debug  bool
}

// WithDB sets the database.
func WithDB(db *freegeoip.DB) Option {
	return func(o *Options) {
		o.DB = db
	}
}

// WithDebug sets the debug flag.
func WithDebug(debug bool) Option {
	return func(o *Options) {
		o.Debug = debug
	}
}

// WithLogger sets the logger.
func WithLogger(logger *zap.Logger) Option {
	return func(o *Options) {
		o.Logger = logger
	}
}
