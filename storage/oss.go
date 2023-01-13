package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

var OSS *OSSClient

type OSSClient struct {
	client   *oss.Client
	bucket   string
	endpoint string
}

func (c *OSSClient) Put(key string, data io.Reader) string {
	bucket, err := c.client.Bucket(c.bucket)
	if err != nil {
		log.Default().Println(err)
	}
	err = bucket.PutObject(key, data)
	if err != nil {
		log.Default().Println(err)
	}
	return "http://" + c.bucket + "." + c.endpoint + "/" + key
}

func NewOSSClient(client *oss.Client, bucket, endpoint string) *OSSClient {
	return &OSSClient{
		client:   client,
		bucket:   bucket,
		endpoint: endpoint,
	}
}
