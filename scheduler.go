package main

import (
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	// "github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"Scheduler/internal/config"
	"Scheduler/internal/server"
	"Scheduler/internal/svc"
	"Scheduler/scheduler"
)

var configFile = flag.String("f", "etc/scheduler.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	if !c.Kubernetes.InCluster && c.Kubernetes.KubeConfig == "" {
		panic("Kubernetes.InCluster or Kubernetes.KubeConfig must be set")
	}
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		scheduler.RegisterSchedulerServer(grpcServer, server.NewSchedulerServer(ctx))
		reflection.Register(grpcServer) // 启用 gRPC reflection
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
