package rpc

import (
	consts "douyin_rpc/server/cmd/api/constant"
	"douyin_rpc/server/cmd/api/global"
	"douyin_rpc/server/cmd/api/kitex_gen/video/feedservice"
	"douyin_rpc/server/cmd/api/middleware"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	nacos "github.com/kitex-contrib/registry-nacos/resolver"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func initVideo() {
	// init resolver
	// Read configuration information from nacos
	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port:   global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         global.ServerConfig.VideoSrvInfo.Namespace,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              consts.NacosLogDir,
		CacheDir:            consts.NacosCacheDir,
		LogLevel:            consts.NacosLogLevel,
	}

	nacosCli, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		})
	r := nacos.NewNacosResolver(nacosCli, nacos.WithGroup(consts.VideoGroup))
	if err != nil {
		klog.Fatalf("new consul client failed: %s", err.Error())
	}
	//// init OpenTelemetry
	//provider.NewOpenTelemetryProvider(
	//	provider.WithServiceName(global.ServerConfig.UserSrvInfo.Name),
	//	provider.WithExportEndpoint(global.ServerConfig.OtelInfo.EndPoint),
	//	provider.WithInsecure(),
	//)

	// create a new client
	c, err := feedservice.NewClient(
		global.ServerConfig.VideoSrvInfo.Name,
		client.WithResolver(r),                                     // service discovery
		client.WithLoadBalancer(loadbalance.NewWeightedBalancer()), // load balance
		client.WithMuxConnection(1),                                // multiplexing
		client.WithMiddleware(middleware.CommonMiddleware),
		client.WithInstanceMW(middleware.ClientMiddleware),
		//client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: global.ServerConfig.VideoSrvInfo.Name}),
	)
	if err != nil {
		klog.Fatalf("ERROR: cannot init client: %v\n", err)
	}
	global.VideoClient = c
}
