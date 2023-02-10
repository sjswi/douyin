package main

import (
	"douyin_rpc/server/cmd/relation/constant"
	"douyin_rpc/server/cmd/relation/initialize"
	relation "douyin_rpc/server/cmd/relation/kitex_gen/relation/relationservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10004)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
	//	provider.WithExportEndpoint("http://192.168.56.102:14268/api/traces"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10004")
	if err != nil {
		return
	}
	svr := relation.NewServer(new(RelationServiceImpl),
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
