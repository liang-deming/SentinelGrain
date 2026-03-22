package client

import (
	"context"

	"SentinelGrain/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CheckRequest  = pb.CheckRequest
	CheckResponse = pb.CheckResponse

	Sentinel interface {
		Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error)
	}

	defaultSentinel struct {
		cli zrpc.Client
	}
)

func NewSentinel(cli zrpc.Client) Sentinel {
	return &defaultSentinel{
		cli: cli,
	}
}

func (m *defaultSentinel) Check(ctx context.Context, in *CheckRequest, opts ...grpc.CallOption) (*CheckResponse, error) {
	c := pb.NewSentinelClient(m.cli.Conn())
	return c.Check(ctx, in, opts...)
}
