package storage

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
)

type OSSClient struct {
	Client   *oss.Client
	Bucket   string
	Endpoint string
}

func (c *OSSClient) Put(key string, data io.Reader) string {
	bucket, err := c.Client.Bucket(c.Bucket)
	if err != nil {
		log.Default().Println(err)
	}
	err = bucket.PutObject(key, data)
	if err != nil {
		log.Default().Println(err)
	}
	return "http://" + c.Bucket + "." + c.Endpoint + "/" + key
}

func NewOSSClient(client *oss.Client, bucket, endpoint string) *OSSClient {
	return &OSSClient{
		Client:   client,
		Bucket:   bucket,
		Endpoint: endpoint,
	}
}
