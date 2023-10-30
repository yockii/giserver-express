package service

import (
	"errors"
	"sync"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/pkg/database"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/util"
)

var StoreService = &storeService{
	cacheByName: make(map[string]*model.Store),
	cacheById:   make(map[database.Int64]*model.Store),
}

type storeService struct {
	cacheByName map[string]*model.Store
	cacheById   map[database.Int64]*model.Store
	lock        sync.Mutex
}

func (s *storeService) Add(store *model.Store) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.Store{
		Name: store.Name,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	store.Id = database.Int64(util.SnowflakeId())

	var omitCols []string
	if store.Name == "" {
		omitCols = append(omitCols, "name")
	}
	if store.StoreType == 0 {
		omitCols = append(omitCols, "store_type")
	}

	_, err = database.DB.Omit(omitCols...).Insert(store)
	if err != nil {
		logger.Errorln(err)
		return
	}
	s.cache(store)
	return
}

func (s *storeService) Update(store *model.Store) error {
	if store.Id == 0 {
		return errors.New("ID must be provided")
	}

	_, err := database.DB.ID(store.Id).Update(store)
	if err != nil {
		logger.Errorln(err)
		return err
	}
	s.cache(store)
	delete(OssService.ossClients, store.Id)
	return nil
}

func (*storeService) Delete(id database.Int64) (err error) {
	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	_, err = sess.Delete(&model.Store{Id: id})
	if err != nil {
		logger.Errorln(err)
		return
	}
	err = sess.Commit()
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (*storeService) List(condition *model.Store, offset, limit int, orderBy string) (total int64, list []*model.Store, err error) {
	if offset < 0 {
		offset = 0
	}
	if limit > 0 {
		if limit > 1000 {
			limit = 1000
		}
	} else {
		limit = 1000
	}
	sess := database.DB.Limit(limit, offset)
	if condition.Name != "" {
		sess.Where("name like ?", "%"+condition.Name+"%")
		condition.Name = ""
	}
	if orderBy != "" {
		sess.OrderBy(orderBy)
	}
	total, err = sess.FindAndCount(&list, condition)
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *storeService) FindByName(name string) (*model.Store, error) {
	store := s.getFromCacheByName(name)
	if store != nil {
		return store, nil
	}
	store = new(model.Store)
	store.Name = name
	exist, err := database.DB.Get(store)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exist {
		s.cache(store)
		return store, nil
	}
	return nil, nil
}

func (s *storeService) GetById(id database.Int64) (*model.Store, error) {
	store := s.getFromCacheById(id)
	if store != nil {
		return store, nil
	}
	store = &model.Store{Id: id}
	exists, err := database.DB.Get(store)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exists {
		s.cache(store)
		return store, nil
	}
	return nil, nil
}

func (s *storeService) getFromCacheByName(storeName string) *model.Store {
	store, ok := s.cacheByName[storeName]
	if !ok {
		return nil
	}
	return store
}
func (s *storeService) getFromCacheById(id database.Int64) *model.Store {
	store, ok := s.cacheById[id]
	if !ok {
		return nil
	}
	return store
}

func (s *storeService) cache(store *model.Store) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.cacheByName[store.Name] = store
	s.cacheById[store.Id] = store
}
