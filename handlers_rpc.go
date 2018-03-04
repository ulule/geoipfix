package ipfix

import (
	"context"

	"github.com/ulule/ipfix/proto"
)

type rpcHandler struct{}

func (h *rpcHandler) GetLocation(context.Context, *proto.GetLocationRequest) (*proto.Location, error) {
	return nil, nil
}
