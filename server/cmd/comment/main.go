package main

import (
	"douyin_rpc/server/cmd/comment/constant"
	"douyin_rpc/server/cmd/comment/initialize"
	"douyin_rpc/server/cmd/comment/initialize/rpc"
	comment "douyin_rpc/server/cmd/comment/kitex_gen/comment/commentservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10006)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	rpc.Init()
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
	//	provider.WithExportEndpoint("http://192.168.56.102:14268/api/traces"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10006")
	if err != nil {
		return
	}
	svr := comment.NewServer(new(CommentServiceImpl),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithServiceAddr(addr),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.DefaultName + "_" + constant.SERVICEName}),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
