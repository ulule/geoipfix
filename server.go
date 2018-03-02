package ipfix

import (
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/valve"
)

// Run launchs the server with a config path.
func Run(config string) error {
	app, err := NewApplication(config)
	if err != nil {
		return err
	}

	defer app.Logger.Sync()

	httpServer := NewHTTPServer(app.Config.Server.HTTP,
		WithLogger(app.Logger),
		WithDebug(app.Config.Debug),
		WithDB(app.DB))
	err = httpServer.Init()
	if err != nil {
		return err
	}

	valv := valve.New()
	baseCtx := valv.Context()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go httpServer.Serve(baseCtx)

	<-c

	app.Logger.Info("shutting down..")

	// first valv
	valv.Shutdown(20 * time.Second)

	httpServer.Shutdown()

	return nil
}
