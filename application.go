package ipfix

import (
	"github.com/fiorix/freegeoip"
	"go.uber.org/zap"
)

// application is the ipfix application.
type application struct {
	DB     *freegeoip.DB
	Logger *zap.Logger
	Config *config
}

// newApplication initializes a new Application instance.
func newApplication(config string) (*application, error) {
	cfg, err := loadConfig(config)
	if err != nil {
		return nil, err
	}

	db, err := openDB(cfg.DatabasePath, UpdateInterval, RetryInterval)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()

	return &application{
		Config: cfg,
		Logger: logger,
		DB:     db,
	}, nil
}
