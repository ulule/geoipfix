package ipfix

import (
	"net/url"
	"time"

	"github.com/fiorix/freegeoip"
)

func openDB(dsn string, updateIntvl time.Duration, maxRetryIntvl time.Duration) (db *freegeoip.DB, err error) {
	u, err := url.Parse(dsn)
	if err != nil || len(u.Scheme) == 0 {
		db, err = freegeoip.Open(dsn)
	} else {
		db, err = freegeoip.OpenURL(dsn, updateIntvl, maxRetryIntvl)
	}
	return
}
