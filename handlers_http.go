package ipfix

import (
	"net"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// Handler is an http.Handler with ipfix options.
type Handler func(opt Options, w http.ResponseWriter, r *http.Request) error

// IPAddressHandler retrieves the IP from request.
func IPAddressHandler(opt Options, w http.ResponseWriter, r *http.Request) error {
	rawIP := chi.URLParam(r, "ipAddress")
	if rawIP == "" {
		rawIP = r.RemoteAddr
	}

	ip := net.ParseIP(rawIP)

	q := &geoipQuery{}
	err := opt.DB.Lookup(ip, &q)
	if err != nil {
		return err
	}

	resp := q.Record(ip, r.Header.Get("Accept-Language"))

	render.JSON(w, r, resp)

	return nil
}
