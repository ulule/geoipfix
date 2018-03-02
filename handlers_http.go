package ipfix

import (
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

// Handler is an http.Handler with ipfix options.
type Handler func(opt Options, w http.ResponseWriter, r *http.Request) error

// IPAddressHandler retrieves the IP from request.
func IPAddressHandler(opt Options, w http.ResponseWriter, r *http.Request) error {
	rawIP := chi.URLParam(r, "ipAddress")
	if rawIP == "" {
		rawIP = r.RemoteAddr
	}

	opt.Logger.Info("Retrieve IP Address from request", zap.String("ip_address", rawIP))

	ip := net.ParseIP(rawIP)
	if ip == nil {
		opt.Logger.Error("IP Address cannot be parsed", zap.String("ip_address", rawIP))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return nil
	}

	q := &geoipQuery{}
	err := opt.DB.Lookup(ip, &q)
	if err != nil {
		return errors.Wrapf(err, "Cannot retrieve geoip information for %s", rawIP)
	}

	resp := q.Record(ip, r.Header.Get("Accept-Language"))

	render.JSON(w, r, resp)

	return nil
}
