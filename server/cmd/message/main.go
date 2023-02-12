package main

import (
	"douyin_rpc/server/cmd/message/constant"
	"douyin_rpc/server/cmd/message/initialize"
	message "douyin_rpc/server/cmd/message/kitex_gen/message/messageservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10003)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
	//	provider.WithExportEndpoint("http://192.168.56.102:14268/api/traces"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10003")
	if err != nil {
		return
	}
	svr := message.NewServer(new(MessageServiceImpl),
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
