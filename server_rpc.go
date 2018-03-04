package ipfix

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/ulule/ipfix/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// rpcServer is an RPC server.
type rpcServer struct {
	srv *grpc.Server
	opt options
	cfg serverRPCConfig
}

func newRPCServer(cfg serverRPCConfig, opts ...option) *rpcServer {
	opt := newOptions(opts...)
	opt.Logger = opt.Logger.With(zap.String("server", "rpc"))

	return &rpcServer{
		cfg: cfg,
		opt: opt,
	}
}

// Init initializes rpc server instance.
func (h *rpcServer) Init() error {
	s := grpc.NewServer()
	proto.RegisterIpfixServer(s, &rpcHandler{opt: h.opt})

	h.srv = s

	return nil
}

// Serve serves rpc requests.
func (h *rpcServer) Serve(ctx context.Context) error {
	addr := fmt.Sprintf(":%s", strconv.Itoa(h.cfg.Port))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	h.opt.Logger.Info("Launch RPC server", zap.String("addr", addr))

	err = h.srv.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown stops the rpc server.
func (h *rpcServer) Shutdown() error {
	h.srv.GracefulStop()

	h.opt.Logger.Info("RPC server shutdown")

	return nil
}
