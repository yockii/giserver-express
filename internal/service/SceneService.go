package service

import (
	"errors"
	"strconv"
	"strings"

	logger "github.com/sirupsen/logrus"

	"github.com/yockii/giserver-express/internal/domain"
	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/util"
)

var SceneService = new(sceneService)

type sceneService struct{}

func (s *sceneService) Add(scene *model.Scene) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.Scene{
		SpaceId: scene.SpaceId,
		Name:    scene.Name,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	scene.Id = util.SnowflakeId()

	var omitCols []string
	if scene.ResourceConfigId == "" {
		omitCols = append(omitCols, "resource_config_id")
	}
	if scene.SupportedMediaTypes == "" {
		omitCols = append(omitCols, "supported_media_types")
	}
	if scene.ResourceType == "" {
		omitCols = append(omitCols, "resource_type")
	}
	if scene.MinCameraDistance == 0 {
		omitCols = append(omitCols, "min_camera_distance")
	}
	if scene.MaxCameraDistance == 0 {
		omitCols = append(omitCols, "max_camera_distance")
	}
	if scene.ScaleLegendVisible == 0 {
		omitCols = append(omitCols, "scale_legend_visible")
	}
	if scene.CameraFov == 0 {
		omitCols = append(omitCols, "camera_fov")
	}
	if scene.FogVisibleAltitude == 0 {
		omitCols = append(omitCols, "fog_visible_altitude")
	}
	if scene.SceneType == "" {
		omitCols = append(omitCols, "scene_type")
	}
	if scene.TerrainExaggeration == 0 {
		omitCols = append(omitCols, "terrain_exaggeration")
	}

	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()
	_, err = sess.Omit(omitCols...).Insert(scene)
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Cols("id", "scene_id").Insert(&model.SceneAtmosphere{
		Id:      util.SnowflakeId(),
		SceneId: scene.Id,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Cols("id", "scene_id").Insert(&model.SceneLatLonGrid{
		Id:      util.SnowflakeId(),
		SceneId: scene.Id,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Cols("id", "scene_id").Insert(&model.SceneCamera{
		Id:      util.SnowflakeId(),
		SceneId: scene.Id,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Cols("id", "scene_id").Insert(&model.SceneFog{
		Id:      util.SnowflakeId(),
		SceneId: scene.Id,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	err = sess.Commit()
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (*sceneService) Update(scene *model.Scene) error {
	if scene.Id == 0 {
		return errors.New("ID must be provided")
	}

	_, err := database.DB.ID(scene.Id).Update(scene)
	if err != nil {
		logger.Errorln(err)
	}
	return err
}

func (*sceneService) Delete(id int64) (err error) {
	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	// 删除子场景下的所有图层
	_, err = sess.Delete(&model.SceneLayer{SceneId: id})
	if err != nil {
		logger.Errorln(err)
		return
	}
	// 删除子场景下所有附加信息
	_, err = sess.Delete(&model.SceneFog{SceneId: id})
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Delete(&model.SceneCamera{SceneId: id})
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Delete(&model.SceneLatLonGrid{SceneId: id})
	if err != nil {
		logger.Errorln(err)
		return
	}
	_, err = sess.Delete(&model.SceneAtmosphere{SceneId: id})
	if err != nil {
		logger.Errorln(err)
		return
	}

	_, err = sess.Delete(&model.Scene{Id: id})
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

func (*sceneService) List(condition *model.Scene, offset, limit int) (total int64, list []*model.Scene, err error) {
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

func (*sceneService) FindByName(name string) (*model.Scene, error) {
	scene := new(model.Scene)
	scene.Name = name
	exist, err := database.DB.Get(scene)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if exist {
		return scene, nil
	}
	return nil, nil
}

func (*sceneService) FindSpaceScenes(spaceId int64) (scenes []*model.Scene, err error) {
	err = database.DB.Cols("id", "name", "resource_config_id", "supported_media_types", "resource_type").Find(&scenes, &model.Scene{SpaceId: spaceId})
	if err != nil {
		logger.Errorln(err)
	}
	return
}

func (s *sceneService) GetRichSceneInfoById(sceneId int64) (*domain.Scene, error) {
	scene := &model.Scene{Id: sceneId}
	exist, err := database.DB.Get(scene)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if !exist {
		return nil, nil
	}
	result := &domain.Scene{
		Name:                scene.Name,
		MinCameraDistance:   scene.MinCameraDistance,
		MaxCameraDistance:   scene.MaxCameraDistance,
		ScaleLegendVisible:  scene.ScaleLegendVisible == 1,
		CameraFOV:           scene.CameraFov,
		FogVisibleAltitude:  scene.FogVisibleAltitude,
		SceneType:           scene.SceneType,
		TerrainExaggeration: scene.TerrainExaggeration,
		Layers:              nil,
	}
	// 大气
	{
		atmosphere := &model.SceneAtmosphere{SceneId: sceneId}
		exist, err = database.DB.Get(atmosphere)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		if exist {
			result.Atmosphere = &domain.SceneAtmosphere{
				Visible: atmosphere.Visible == 1,
			}
		}
	}
	// 经纬
	{
		latLonGrid := &model.SceneLatLonGrid{SceneId: sceneId}
		exist, err = database.DB.Get(latLonGrid)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		if exist {
			result.LatLonGrid = &domain.SceneLatLonGrid{
				Visible:     latLonGrid.Visible == 1,
				TextVisible: latLonGrid.TextVisible == 1,
			}
		}
	}
	// 相机
	{
		camera := &model.SceneCamera{SceneId: sceneId}
		exist, err = database.DB.Get(camera)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		if exist {
			result.Camera = &domain.SceneCamera{
				Altitude:     camera.Altitude,
				Latitude:     camera.Latitude,
				Longitude:    camera.Longitude,
				Heading:      camera.Heading,
				AltitudeMode: camera.AltitudeMode,
				Tilt:         camera.Tilt,
				Empty:        camera.Empty,
			}
		}
	}
	// 雾
	{
		fog := &model.SceneFog{SceneId: sceneId}
		exist, err = database.DB.Get(fog)
		if err != nil {
			logger.Errorln(err)
			return nil, err
		}
		if exist {
			color := &domain.Color{
				Red: 255, Green: 255, Blue: 255, Alpha: 255,
			}
			colorStrs := strings.Split(fog.Color, ",")
			if len(colorStrs) == 4 {
				color.Red, _ = strconv.Atoi(colorStrs[0])
				color.Green, _ = strconv.Atoi(colorStrs[1])
				color.Blue, _ = strconv.Atoi(colorStrs[2])
				color.Alpha, _ = strconv.Atoi(colorStrs[3])
			}
			result.Fog = &domain.SceneFog{
				Mode:          fog.Mode,
				EndDistance:   fog.EndDistance,
				StartDistance: fog.StartDistance,
				Color:         color,
				Density:       fog.Density,
				Enable:        fog.Enable == 1,
			}
		}
	}

	// 层信息数据
	var layers []*domain.SceneLayer
	layers, err = LayerService.FindSceneDomainLayers(sceneId)
	if err != nil {
		return nil, err
	}
	result.Layers = layers
	return result, nil
}

func (*sceneService) AddSceneAtmosphere(atmosphere *model.SceneAtmosphere) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.SceneAtmosphere{
		SceneId: atmosphere.SceneId,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	atmosphere.Id = util.SnowflakeId()
	_, err = database.DB.Insert(atmosphere)
	if err != nil {
		logger.Errorln(err)
	}
	return
}
func (*sceneService) AddSceneLatLonGrid(latLonGrid *model.SceneLatLonGrid) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.SceneLatLonGrid{
		SceneId: latLonGrid.SceneId,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	latLonGrid.Id = util.SnowflakeId()
	_, err = database.DB.Insert(latLonGrid)
	if err != nil {
		logger.Errorln(err)
	}
	return
}
func (*sceneService) AddSceneCamera(camera *model.SceneCamera) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.SceneCamera{
		SceneId: camera.SceneId,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	camera.Id = util.SnowflakeId()
	_, err = database.DB.Insert(camera)
	if err != nil {
		logger.Errorln(err)
	}
	return
}
func (*sceneService) AddSceneFog(fog *model.SceneFog) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.SceneFog{
		SceneId: fog.SceneId,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}
	fog.Id = util.SnowflakeId()
	_, err = database.DB.Insert(fog)
	if err != nil {
		logger.Errorln(err)
	}
	return
}
