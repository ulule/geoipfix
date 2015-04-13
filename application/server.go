package application

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/fiorix/freegeoip"
	"github.com/rs/cors"
	"github.com/thoas/stats"
	"github.com/tylerb/graceful"
	"net/http"
	"net/url"
	"runtime"
	"strconv"
	"time"
)

func Run(config string) error {
	runtime.GOMAXPROCS(runtime.NumCPU())

	app, err := NewFromConfigPath(config)

	if err != nil {
		return err
	}

	databasePath, err := app.Jq.String("database_path")

	if err != nil {
		databasePath = DatabaseURL
	}

	db, err := openDB(databasePath, UpdateInterval, RetryInterval)

	if err != nil {
		return err
	}

	mux := http.NewServeMux()

	middleware := stats.New()

	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		stats := middleware.Data()

		b, _ := json.Marshal(stats)

		w.Write(b)
	})

	mux.Handle("/json/", freegeoip.ProxyHandler(freegeoip.NewHandler(db, &freegeoip.JSONEncoder{})))

	allowedOrigins, _ := app.Jq.ArrayOfStrings("allowed_origins")
	allowedMethods, _ := app.Jq.ArrayOfStrings("allowed_methods")

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(mux)
	n.Use(cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: allowedMethods,
	}))
	graceful.Run(fmt.Sprintf(":%s", strconv.Itoa(app.Port())), 10*time.Second, n)

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
