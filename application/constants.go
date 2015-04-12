package application

import (
	"time"
)

const Version = "0.1"
const DefaultPort = 3001
const DatabaseURL = "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
const UpdateInterval = 24 * time.Hour
const RetryInterval = time.Hour
