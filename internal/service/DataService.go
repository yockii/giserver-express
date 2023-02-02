package service

import (
	"errors"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/util"
)

var DataService = &dataService{
	cache: make(map[string]map[string]*model.Data),
}

type dataService struct {
	cache map[string]map[string]*model.Data
}

func (*dataService) Add(data *model.Data) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.Data{
		SpaceId: data.SpaceId,
		Name:    data.Name,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	data.Id = util.SnowflakeId()

	var omitCols []string
	if data.Name == "" {
		omitCols = append(omitCols, "name")
	}
	if data.DataType == "" {
		omitCols = append(omitCols, "data_type")
	}

	_, err = database.DB.Omit(omitCols...).Insert(data)
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (*dataService) Update(data *model.Data) error {
	if data.Id == 0 {
		return errors.New("ID must be provided")
	}

	_, err := database.DB.ID(data.Id).Update(data)
	if err != nil {
		logger.Errorln(err)
	}
	return err
}

func (*dataService) Delete(id int64) (err error) {
	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	_, err = sess.Delete(&model.Data{Id: id})
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

func (*dataService) List(condition *model.Data, offset, limit int, orderBy string) (total int64, list []*model.Data, err error) {
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

func (*dataService) FindByName(name string) (*model.Data, error) {
	data := new(model.Data)
	data.Name = name
	exist, err := database.DB.Get(data)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exist {
		return data, nil
	}
	return nil, nil
}

func (*dataService) FindSpaceDataList(spaceId int64) (datas []*model.Data, err error) {
	err = database.DB.Find(&datas, &model.Data{SpaceId: spaceId})
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *dataService) GetBySpaceIdAndDataName(spaceId int64, name string) (*model.Data, error) {
	data := &model.Data{
		SpaceId: spaceId,
		Name:    name,
	}
	exists, err := database.DB.Get(data)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exists {
		return data, nil
	}
	return nil, nil
}

func (s *dataService) GetById(id int64) (*model.Data, error) {
	data := &model.Data{Id: id}

	exists, err := database.DB.Get(data)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exists {
		return data, nil
	}
	return nil, nil
}

func (s *dataService) GetFromCache(spaceName string, dataName string) *model.Data {
	spaceDataMap, ok := s.cache[spaceName]
	if !ok {
		return nil
	}
	data, has := spaceDataMap[dataName]
	if !has {
		return nil
	}
	return data
}

func (s *dataService) Cache(spaceName string, data *model.Data) {
	spaceDataMap, ok := s.cache[spaceName]
	if !ok || spaceDataMap == nil {
		spaceDataMap = make(map[string]*model.Data)
		s.cache[spaceName] = spaceDataMap
	}
	spaceDataMap[data.Name] = data
}
