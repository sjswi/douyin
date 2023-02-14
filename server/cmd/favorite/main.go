package main

import (
	"context"
	"douyin_rpc/server/cmd/favorite/constant"
	"douyin_rpc/server/cmd/favorite/global"
	"douyin_rpc/server/cmd/favorite/initialize"
	"douyin_rpc/server/cmd/favorite/initialize/rpc"
	favorite "douyin_rpc/server/cmd/favorite/kitex_gen/favorite/favoriteservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10002)
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
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10002")
	if err != nil {
		return
	}
	svr := favorite.NewServer(new(FavoriteServiceImpl),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithRegistryInfo(info),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.DefaultName + "_" + constant.SERVICEName}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
