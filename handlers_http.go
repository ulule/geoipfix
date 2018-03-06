package ipfix

import (
	"net"
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
)

type httpHandler struct {
	options
}

// GetLocation retrieves the IP from request.
func (h *httpHandler) GetLocation(w http.ResponseWriter, r *http.Request) error {
	rawIP := chi.URLParam(r, "ipAddress")
	if rawIP == "" {
		rawIP = r.RemoteAddr
	}

	log := h.Logger.With(zap.String("ip_address", rawIP))
	log.Info("Retrieve IP Address from request", zap.String("ip_address", rawIP))

	ip := net.ParseIP(rawIP)
	if ip == nil {
		log.Error("IP Address cannot be parsed", zap.String("ip_address", rawIP))

		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)

		return nil
	}

	q := geoipQuery{}
	err := h.DB.Lookup(ip, &q)
	if err != nil {
		return errors.Wrapf(err, "Cannot retrieve geoip information for %s", rawIP)
	}

	resp := q.Record(ip, r.Header.Get("Accept-Language"))

	render.JSON(w, r, resp)

	return nil
}
