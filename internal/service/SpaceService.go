package service

import (
	"errors"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/util"
)

var SpaceService = new(spaceService)

type spaceService struct{}

func (*spaceService) Add(space *model.Space) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.Space{
		Name: space.Name,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	space.Id = util.SnowflakeId()
	_, err = database.DB.Insert(space)
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (*spaceService) Update(space *model.Space) error {
	if space.Id == 0 {
		return errors.New("ID must be provided")
	}

	_, err := database.DB.ID(space.Id).Update(space)
	if err != nil {
		logger.Errorln(err)
	}
	return err
}

func (*spaceService) Delete(id int64) (err error) {
	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	// 所有的子场景
	var scenes []*model.Scene
	if err = sess.Find(&scenes, &model.Scene{SpaceId: id}); err != nil {
		logger.Errorln(err)
		return
	}
	for _, scene := range scenes {
		s := scene
		// 删除子场景下的所有图层
		_, err = sess.Delete(&model.SceneLayer{SceneId: s.Id})
		if err != nil {
			logger.Errorln(err)
			return
		}
		// 删除子场景下所有附加信息
		_, err = sess.Delete(&model.SceneFog{SceneId: s.Id})
		if err != nil {
			logger.Errorln(err)
			return
		}
		_, err = sess.Delete(&model.SceneCamera{SceneId: s.Id})
		if err != nil {
			logger.Errorln(err)
			return
		}
		_, err = sess.Delete(&model.SceneLatLonGrid{SceneId: s.Id})
		if err != nil {
			logger.Errorln(err)
			return
		}
		_, err = sess.Delete(&model.SceneAtmosphere{SceneId: s.Id})
		if err != nil {
			logger.Errorln(err)
			return
		}
	}
	// 删除所有子场景
	_, err = sess.Delete(&model.Scene{SpaceId: id})
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 删除空间
	_, err = sess.Delete(&model.Space{Id: id})
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

func (*spaceService) List(condition *model.Space, offset, limit int) (total int64, list []*model.Space, err error) {
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

func (*spaceService) FindByName(name string) (*model.Space, error) {
	space := new(model.Space)
	space.Name = name
	exist, err := database.DB.Get(space)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exist {
		return space, nil
	}
	return nil, nil
}
