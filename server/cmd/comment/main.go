package main

import (
	"context"
	"douyin_rpc/server/cmd/comment/constant"
	"douyin_rpc/server/cmd/comment/global"
	"douyin_rpc/server/cmd/comment/initialize"
	"douyin_rpc/server/cmd/comment/initialize/rpc"
	comment "douyin_rpc/server/cmd/comment/kitex_gen/comment/commentservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10001)
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
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10001")
	if err != nil {
		log.Println(err.Error())
		return
	}
	svr := comment.NewServer(new(CommentServiceImpl),
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
