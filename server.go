package ipfix

import (
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/valve"
)

// Run launchs the server with a config path.
func Run(config string) error {
	app, err := newApplication(config)
	if err != nil {
		return err
	}

	defer app.Logger.Sync()

	httpServer := newHTTPServer(app.Config.Server.HTTP,
		withLogger(app.Logger),
		withDebug(app.Config.Debug),
		withDB(app.DB))
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

	app.Logger.Info("Shutting down servers..")

	// first valv
	valv.Shutdown(20 * time.Second)

	return httpServer.Shutdown()
}
