package service

import (
	"io"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/aliyun"
)

var OssService = &ossService{ossClients: make(map[int64]*aliyun.OssClient)}

type ossService struct {
	ossClients map[int64]*aliyun.OssClient
}

func (s *ossService) StreamFromStore(store *model.Store, objectKey string) (reader io.Reader, err error) {
	client, ok := s.ossClients[store.Id]
	if !ok {
		client, err = aliyun.NewOssClient(store.Endpoint, store.AccessKeyId, store.AccessKeySecret, store.BucketName, store.UseCname == 1)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		s.ossClients[store.Id] = client
	}
	return client.ReadStream(objectKey)
}
