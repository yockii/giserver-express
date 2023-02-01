package aliyun

import (
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	logger "github.com/sirupsen/logrus"
)

type OssClient struct {
	Client *oss.Client
	Bucket *oss.Bucket
}

func NewOssClient(endpoint, bucketName, accessKeyId, accessKeySecret string, selfDomain bool) (*OssClient, error) {
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret, oss.UseCname(selfDomain))
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	var bucket *oss.Bucket
	bucket, err = client.Bucket(bucketName)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}

	return &OssClient{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (c *OssClient) ReadStream(objectKey string) (io.ReadCloser, error) {
	body, err := c.Bucket.GetObject(objectKey)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	return body, nil
}
