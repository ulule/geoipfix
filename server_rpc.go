package ipfix

import (
	"context"
	"fmt"
	"net"
	"strconv"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/ulule/ipfix/proto"
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
	grpc_zap.ReplaceGrpcLogger(h.opt.Logger)

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_zap.UnaryServerInterceptor(h.opt.Logger),
		),
	)
	proto.RegisterIpfixServer(s, &rpcHandler{h.opt})

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
