package main

import (
	"context"
	"douyin_rpc/server/cmd/user/constant"
	"douyin_rpc/server/cmd/user/global"
	"douyin_rpc/server/cmd/user/initialize"
	"douyin_rpc/server/cmd/user/initialize/rpc"
	user "douyin_rpc/server/cmd/user/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10005)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	rpc.Init()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
		provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10005")
	if err != nil {
		return
	}
	svr := user.NewServer(new(UserServiceImpl),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.DefaultName + "_" + constant.SERVICEName}),
		//server.WithTracer(prometheus.NewServerTracer(":9092", "/userServer")),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
