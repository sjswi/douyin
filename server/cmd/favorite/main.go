package main

import (
	"douyin_rpc/server/cmd/favorite/constant"
	"douyin_rpc/server/cmd/favorite/initialize"
	"douyin_rpc/server/cmd/favorite/initialize/rpc"
	favorite "douyin_rpc/server/cmd/favorite/kitex_gen/favorite/favoriteservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10005)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	rpc.Init()
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
	//	provider.WithExportEndpoint("192.168.56.102:14268"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10005")
	if err != nil {
		return
	}
	svr := favorite.NewServer(new(FavoriteServiceImpl),
		server.WithRegistry(r),
		server.WithServiceAddr(addr),
		server.WithRegistryInfo(info),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.DefaultName + "_" + constant.SERVICEName}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
