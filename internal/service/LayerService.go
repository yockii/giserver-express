package service

import (
	"errors"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/domain"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/util"
)

var LayerService = new(layerService)

type layerService struct{}

func (*layerService) Add(layer *model.SceneLayer) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.SceneLayer{
		SceneId: layer.SceneId,
		DataId:  layer.DataId,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	data := &model.Data{Id: layer.DataId}
	exists, err := database.DB.Get(data)
	if err != nil {
		logger.Errorln(err)
		return
	}
	if !exists {
		logger.Debugln("data not exists...")
		return
	}

	logger.Infoln("ready to add layer...")

	layer.Name = data.DataName
	layer.Id = util.SnowflakeId()

	var omitCols []string
	if layer.MinVisibleAltitude == 0 {
		omitCols = append(omitCols, "min_visible_altitude")
	}
	if layer.MaxVisibleAltitude == 0 {
		omitCols = append(omitCols, "max_visible_altitude")
	}
	if layer.VisibleDistance == 0 {
		omitCols = append(omitCols, "visible_distance")
	}
	if layer.IsWebDatasource == 0 {
		omitCols = append(omitCols, "is_web_datasource")
	}
	if layer.AlwaysRender == 0 {
		omitCols = append(omitCols, "always_render")
	}
	if layer.Visible == 0 {
		omitCols = append(omitCols, "visible")
	}
	if layer.Level == 0 {
		omitCols = append(omitCols, "Level")
	}
	if layer.UseTwoDimenCache == 0 {
		omitCols = append(omitCols, "use_two_dimen_Cache")
	}
	if layer.Editable == 0 {
		omitCols = append(omitCols, "editable")
	}
	if layer.Caption == "" {
		omitCols = append(omitCols, "caption")
	}
	if layer.HasLocalCache == 0 {
		omitCols = append(omitCols, "has_local_cache")
	}
	if layer.Layer3DType == "" {
		omitCols = append(omitCols, "layer_3d_type")
	}

	_, err = database.DB.Omit(omitCols...).Insert(layer)
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (*layerService) Update(layer *model.SceneLayer) error {
	if layer.Id == 0 {
		return errors.New("ID must be provided")
	}

	_, err := database.DB.ID(layer.Id).Update(layer)
	if err != nil {
		logger.Errorln(err)
	}
	return err
}

func (*layerService) Delete(id int64) (err error) {
	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	_, err = sess.Delete(&model.SceneLayer{Id: id})
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

func (*layerService) List(condition *model.SceneLayer, offset, limit int, orderBy string) (total int64, list []*model.SceneLayer, err error) {
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

func (*layerService) FindByName(name string) (*model.SceneLayer, error) {
	layer := new(model.SceneLayer)
	layer.Name = name
	exist, err := database.DB.Get(layer)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exist {
		return layer, nil
	}
	return nil, nil
}

func (*layerService) FindSceneLayers(sceneId int64) (layers []*model.SceneLayer, err error) {
	err = database.DB.Find(&layers, &model.SceneLayer{SceneId: sceneId})
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (*layerService) FindSceneDomainLayers(sceneId int64) (list []*domain.SceneLayer, err error) {
	var layers []*model.SceneLayer
	err = database.DB.Find(&layers, &model.SceneLayer{SceneId: sceneId})
	if err != nil {
		logger.Errorln(err)
	}

	for _, layer := range layers {
		l := layer
		data := &model.Data{Id: l.DataId}
		exists, _ := database.DB.Get(data)
		if exists {
			dl := &domain.SceneLayer{
				Name:               l.Name,
				MinVisibleAltitude: l.MinVisibleAltitude,
				MaxVisibleAltitude: l.MaxVisibleAltitude,
				VisibleDistance:    l.VisibleDistance,
				IsWebDatasource:    l.IsWebDatasource == 1,
				AlwaysRender:       l.AlwaysRender == 1,
				Visible:            l.Visible == 1,
				Level:              l.Level,
				UseTwoDimenCache:   l.UseTwoDimenCache == 1,
				Editable:           l.Editable == 1,
				Caption:            l.Caption,
				Description:        l.Description,
				DataName:           data.DataName,
				HasLocalCache:      l.HasLocalCache == 1,
				Layer3DType:        l.Layer3DType,
				DataConfigPath:     data.DataConfigPath,
				ExtendXML:          nil,
				SubLayers:          nil,
				Type:               nil,
				ParentLayerName:    nil,
				Bounds:             nil,
				CachePassword:      nil,
				Queryable:          false,
				OldCache:           false,
			}
			list = append(list, dl)
		}
	}

	return list, nil
}

func (s *layerService) GetBySceneIdAndLayerName(sceneId int64, name string) (*model.SceneLayer, error) {
	layer := &model.SceneLayer{
		SceneId: sceneId,
		Name:    name,
	}
	exists, err := database.DB.Get(layer)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exists {
		return layer, nil
	}
	return nil, nil
}

func (s *layerService) GetBySpaceIdAndLayerName(spaceId int64, layerName string) (*model.SceneLayer, error) {
	layer := &model.SceneLayer{
		SpaceId: spaceId,
		Name:    layerName,
	}
	exists, err := database.DB.Get(layer)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exists {
		return layer, nil
	}
	return nil, nil
}
