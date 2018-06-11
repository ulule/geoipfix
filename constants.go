package geoipfix

import (
	"time"
)

// Version is the current application version
const Version = "0.1.0"

// DefaultPort is the default server port
const DefaultPort = 3001

// DatabaseURL is the full url to download the maxmind database
const DatabaseURL = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"

// UpdateInterval is the default time to update the database
const UpdateInterval = 24 * time.Hour

// RetryInterval is the default retry time to retry the update
const RetryInterval = time.Hour

// compilation variables.
var (
	Branch    string
	Revision  string
	BuildTime string
	Compiler  string
)
