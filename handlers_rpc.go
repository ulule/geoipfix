package ipfix

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"github.com/ulule/ipfix/proto"
	"go.uber.org/zap"
)

type rpcHandler struct {
	opt options
}

// GetLocation retrieves location from protobuf
func (h *rpcHandler) GetLocation(ctx context.Context, req *proto.GetLocationRequest) (*proto.Location, error) {
	rawIP := req.IpAddress

	h.opt.Logger.Info("Retrieve IP Address from request", zap.String("ip_address", rawIP))

	ip := net.ParseIP(rawIP)
	if ip == nil {
		h.opt.Logger.Error("IP Address cannot be parsed", zap.String("ip_address", rawIP))

		return nil, errors.Errorf("IP Address %s cannot be parsed", rawIP)
	}

	q := geoipQuery{}
	err := h.opt.DB.Lookup(ip, &q)
	if err != nil {
		return nil, err
	}

	return recordToProto(q.Record(ip, req.Language)), nil
}
