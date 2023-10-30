package service

import (
	"io"
	"sync"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/aliyun"
	"github.com/yockii/giserver-express/pkg/database"
)

var OssService = &ossService{ossClients: make(map[database.Int64]*aliyun.OssClient)}

type ossService struct {
	ossClients map[database.Int64]*aliyun.OssClient
	lock       sync.Mutex
}

func (s *ossService) StreamFromStore(store *model.Store, objectKey string) (reader io.Reader, err error) {
	client, ok := s.ossClients[store.Id]
	if !ok {
		s.lock.Lock()
		defer s.lock.Unlock()

		logger.Debugf("store info: %+v", store)
		client, err = aliyun.NewOssClient(store.Endpoint, store.BucketName, store.AccessKeyId, store.AccessKeySecret, store.UseCname == 1)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		s.ossClients[store.Id] = client
	}
	return client.ReadStream(objectKey)
}
