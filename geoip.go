// This code has been taken from freegeoip repository
// it can be found here: https://github.com/fiorix/freegeoip/blob/master/apiserver/api.go

package geoipfix

import (
	"encoding/xml"
	"math"
	"net"
	"net/url"
	"time"

	"golang.org/x/text/language"

	"github.com/fiorix/freegeoip"
)

type geoipQuery struct {
	freegeoip.DefaultQuery
}

func parseAcceptLanguage(header string, dbLangs map[string]string) string {
	// supported languages -- i.e. languages available in the DB
	matchLangs := []language.Tag{
		language.English,
	}

	// parse available DB languages and add to supported
	for name := range dbLangs {
		matchLangs = append(matchLangs, language.Raw.Make(name))

	}

	var matcher = language.NewMatcher(matchLangs)

	// parse header
	t, _, _ := language.ParseAcceptLanguage(header)
	// match most acceptable language
	tag, _, _ := matcher.Match(t...)
	// extract base language
	base, _ := tag.Base()

	return base.String()

}

func roundFloat(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)

	} else {
		round = math.Floor(digit)

	}
	return round / pow
}

type record struct {
	XMLName     xml.Name `xml:"Response" json:"-"`
	IP          string   `json:"ip"`
	CountryCode string   `json:"country_code"`
	CountryName string   `json:"country_name"`
	RegionCode  string   `json:"region_code"`
	RegionName  string   `json:"region_name"`
	City        string   `json:"city"`
	ZipCode     string   `json:"zip_code"`
	TimeZone    string   `json:"time_zone"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	MetroCode   uint     `json:"metro_code"`
}

func (q *geoipQuery) Record(ip net.IP, lang string) *record {
	lang = parseAcceptLanguage(lang, q.Country.Names)

	r := &record{
		IP:          ip.String(),
		CountryCode: q.Country.ISOCode,
		CountryName: q.Country.Names[lang],
		City:        q.City.Names[lang],
		ZipCode:     q.Postal.Code,
		TimeZone:    q.Location.TimeZone,
		Latitude:    roundFloat(q.Location.Latitude, .5, 4),
		Longitude:   roundFloat(q.Location.Longitude, .5, 4),
		MetroCode:   q.Location.MetroCode,
	}
	if len(q.Region) > 0 {
		r.RegionCode = q.Region[0].ISOCode
		r.RegionName = q.Region[0].Names[lang]

	}
	return r
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
