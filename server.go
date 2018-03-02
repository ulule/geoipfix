package ipfix

// Run launchs the server with a config path.
func Run(config string) error {
	app, err := NewApplication(config)
	if err != nil {
		return err
	}

	httpServer := NewHTTPServer(app.Config.Server.HTTP,
		WithLogger(app.Logger),
		WithDB(app.DB))
	err = httpServer.Init()
	if err != nil {
		return err
	}

	return httpServer.Serve()
}
