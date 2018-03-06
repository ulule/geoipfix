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

	services := []service{}

	if app.Config.Server.HTTP != nil {
		httpServer := newHTTPServer(*app.Config.Server.HTTP,
			withLogger(app.Logger),
			withDebug(app.Config.Debug),
			withDB(app.DB))
		services = append(services, httpServer)
	}

	if app.Config.Server.RPC != nil {
		rpcServer := newRPCServer(*app.Config.Server.RPC,
			withLogger(app.Logger),
			withDB(app.DB))
		services = append(services, rpcServer)
	}

	valv := valve.New()
	baseCtx := valv.Context()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	for _, service := range services {
		err = service.Init()
		if err != nil {
			return err
		}

		go service.Serve(baseCtx)
	}

	<-c

	app.Logger.Info("Shutting down servers..")

	// first valv
	valv.Shutdown(20 * time.Second)

	for _, service := range services {
		err = service.Shutdown()
		if err != nil {
			return nil
		}
	}

	return nil
}
