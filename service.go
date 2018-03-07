package ipfix

import "context"

type service interface {
	Init() error
	Serve(ctx context.Context) error
	Shutdown() error
	Name() string
}
