package main

import (
	"douyin_rpc/server/cmd/user/constant"
	"douyin_rpc/server/cmd/user/initialize"
	user "douyin_rpc/server/cmd/user/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"log"
	"net"
)

func main() {
	r, info := initialize.InitNacos(10002)
	initialize.InitDB()
	initialize.InitRedis()
	initialize.InitLogger()
	//p := provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(constant.DefaultName+"_"+constant.SERVICEName),
	//	provider.WithExportEndpoint("192.168.56.102:4317"),
	//	provider.WithInsecure(),
	//)
	//defer p.Shutdown(context.Background())
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:10002")
	if err != nil {
		return
	}
	svr := user.NewServer(new(UserServiceImpl),
		server.WithRegistry(r),
		server.WithRegistryInfo(info),
		server.WithServiceAddr(addr),
		//server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: constant.DefaultName + "_" + constant.SERVICEName}),
		//server.WithTracer(prometheus.NewServerTracer(":9092", "/userServer")),
	)

	err = svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
