package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/fiorix/freegeoip"
	"github.com/rs/cors"
	"github.com/thoas/stats"
	"github.com/tylerb/graceful"
)

// Run launchs the server with a config path
func Run(config string) error {
	fmt.Println("Loading config from", config, "...")

	jq, err := Load(config)

	if err != nil {
		return err
	}

	databasePath, err := jq.String("database_path")

	if err != nil {
		databasePath = DatabaseURL
	}

	db, err := openDB(databasePath, UpdateInterval, RetryInterval)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	s := stats.New()

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := s.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
	})

	mux.Handle("/json/", freegeoip.ProxyHandler(freegeoip.NewHandler(db, &freegeoip.JSONEncoder{})))

	allowedOrigins, _ := jq.ArrayOfStrings("allowed_origins")
	allowedMethods, _ := jq.ArrayOfStrings("allowed_methods")

	n := negroni.New(negroni.NewRecovery(), NewLogrusMiddleware(), s)
	n.UseHandler(mux)
	n.Use(cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
	}))

	port, err := jq.Int("port")

	if err != nil {
		port = DefaultPort
	}

	fmt.Println("HTTP Server is listening", port)

	graceful.Run(fmt.Sprintf(":%s", strconv.Itoa(port)), 10*time.Second, n)

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
