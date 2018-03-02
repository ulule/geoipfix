package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/fiorix/freegeoip"
	"github.com/rs/cors"
	"github.com/tylerb/graceful"
)

// Run launchs the server with a config path
func Run(config string) error {
	fmt.Println("Loading config from", config, "...")

	cfg, err := Load(config)
	if err != nil {
		return err
	}

	db, err := openDB(cfg.DatabasePath, UpdateInterval, RetryInterval)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	mux.Handle("/json/", freegeoip.ProxyHandler(freegeoip.NewHandler(db, &freegeoip.JSONEncoder{})))

	n := negroni.New(negroni.NewRecovery(), NewLogrusMiddleware())
	n.UseHandler(mux)
	n.Use(cors.New(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: cfg.AllowedMethods,
	}))

	fmt.Println("HTTP Server is listening", cfg.Port)

	graceful.Run(fmt.Sprintf(":%s", strconv.Itoa(cfg.Port)), 10*time.Second, n)

	return nil
}

func openDB(dsn string, updateIntvl time.Duration, maxRetryIntvl time.Duration) (db *freegeoip.DB, err error) {
	u, err := url.Parse(dsn)
	if err != nil || len(u.Scheme) == 0 {
		db, err = freegeoip.Open(dsn)
	} else {
		db, err = freegeoip.OpenURL(dsn, updateIntvl, maxRetryIntvl)
	}
	return
}
