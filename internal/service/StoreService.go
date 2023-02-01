package service

import (
	"errors"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/pkg/database"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/util"
)

var StoreService = &storeService{
	cacheByName: make(map[string]*model.Store),
	cacheById:   make(map[int64]*model.Store),
}

type storeService struct {
	cacheByName map[string]*model.Store
	cacheById   map[int64]*model.Store
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
	store.Id = util.SnowflakeId()

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
	return nil
}

func (*storeService) Delete(id int64) (err error) {
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

func (*storeService) List(condition *model.Store, offset, limit int) (total int64, list []*model.Store, err error) {
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

func (s *storeService) GetById(id int64) (*model.Store, error) {
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
func (s *storeService) getFromCacheById(id int64) *model.Store {
	store, ok := s.cacheById[id]
	if !ok {
		return nil
	}
	return store
}

func (s *storeService) cache(store *model.Store) {
	s.cacheByName[store.Name] = store
	s.cacheById[store.Id] = store
}
