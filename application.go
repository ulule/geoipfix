package ipfix

import (
	"github.com/fiorix/freegeoip"
	"go.uber.org/zap"
)

// Application is the ipfix application.
type Application struct {
	DB     *freegeoip.DB
	Logger *zap.Logger
	Config *Config
}

// NewApplication initializes a new Application instance.
func NewApplication(config string) (*Application, error) {
	cfg, err := Load(config)
	if err != nil {
		return nil, err
	}

	db, err := openDB(cfg.DatabasePath, UpdateInterval, RetryInterval)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()

	return &Application{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}, nil
}
