package initialize

import (
	"douyin_rpc/server/cmd/video/global"
	"douyin_rpc/server/cmd/video/storage"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func InitOSS() {

	client, err := oss.New(global.ServerConfig.OSSInfo.Endpoint, global.ServerConfig.OSSInfo.AccessKeyId, global.ServerConfig.OSSInfo.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	global.OSS = storage.NewOSSClient(client, global.ServerConfig.OSSInfo.Bucket, global.ServerConfig.OSSInfo.Endpoint)
}
