package main

import (
	"flag"
	"fmt"
	"time"

	"SentinelGrain/rpc/internal/config"
	"SentinelGrain/rpc/internal/server"
	"SentinelGrain/rpc/internal/svc"
	"SentinelGrain/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Default expects current working directory to be rpc/ (e.g. cd rpc && go run ./cmd/sentinel-rpc).
// From repository root: go run ./rpc/cmd/sentinel-rpc -f rpc/etc/sentinel.yaml
var configFile = flag.String("f", "etc/sentinel.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	intervalSec := c.QuotaRefreshInterval
	if intervalSec <= 0 {
		intervalSec = 5
	}
	ctx.QuotaRules.StartPeriodicRefresh(time.Duration(intervalSec) * time.Second)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSentinelServer(grpcServer, server.NewSentinelServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
