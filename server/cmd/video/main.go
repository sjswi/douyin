package main

import (
	"context"
	"douyin_rpc/server/cmd/video/constant"
	"douyin_rpc/server/cmd/video/global"
	"douyin_rpc/server/cmd/video/initialize"
	"douyin_rpc/server/cmd/video/initialize/rpc"
	video "douyin_rpc/server/cmd/video/kitex_gen/video/feedservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10006)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	initialize.InitOSS()
	rpc.Init()
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
		provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10006")
	if err != nil {
		return
	}
	svr := video.NewServer(new(FeedServiceImpl),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.DefaultName + "_" + constant.SERVICEName}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
