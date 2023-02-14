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

func (s *sceneService) Add(sceneDomain *domain.Scene) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.Scene{
		SpaceId: sceneDomain.SpaceId,
		Name:    sceneDomain.Name,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	scene := sceneDomain.Scene
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
	if !sceneDomain.ScaleLegendVisible || scene.ScaleLegendVisible == 0 {
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
	{
		cols := []string{"id", "scene_id"}
		visible := 1
		if sceneDomain.Atmosphere != nil {
			if !sceneDomain.Atmosphere.Visible {
				visible = -1
				cols = append(cols, "visible")
			}
		}
		_, err = sess.Cols(cols...).Insert(&model.SceneAtmosphere{
			Id:      util.SnowflakeId(),
			SceneId: scene.Id,
			Visible: visible,
		})
		if err != nil {
			logger.Errorln(err)
			return
		}
	}
	{
		cols := []string{"id", "scene_id"}
		latLonGrid := &model.SceneLatLonGrid{
			Id:      util.SnowflakeId(),
			SceneId: scene.Id,
		}
		if sceneDomain.LatLonGrid != nil {
			if sceneDomain.LatLonGrid.Visible {
				latLonGrid.Visible = 1
				cols = append(cols, "visible")
			}
			if !sceneDomain.LatLonGrid.TextVisible {
				latLonGrid.TextVisible = -1
				cols = append(cols, "text_visible")
			}
		}
		_, err = sess.Cols(cols...).Insert(latLonGrid)
		if err != nil {
			logger.Errorln(err)
			return
		}
	}
	{
		cols := []string{"id", "scene_id"}
		camera := &model.SceneCamera{
			Id:      util.SnowflakeId(),
			SceneId: scene.Id,
		}
		if sceneDomain.Camera != nil {
			if sceneDomain.Camera.Altitude != 0 {
				camera.Altitude = sceneDomain.Camera.Altitude
				cols = append(cols, "altitude")
			}
			if sceneDomain.Camera.Latitude != 0 {
				camera.Latitude = sceneDomain.Camera.Latitude
				cols = append(cols, "latitude")
			}
			if sceneDomain.Camera.Longitude != 0 {
				camera.Longitude = sceneDomain.Camera.Longitude
				cols = append(cols, "longitude")
			}
			if sceneDomain.Camera.Heading != 0 {
				camera.Heading = sceneDomain.Camera.Heading
				cols = append(cols, "heading")
			}
			if sceneDomain.Camera.AltitudeMode != "" {
				camera.AltitudeMode = sceneDomain.Camera.AltitudeMode
				cols = append(cols, "altitudeMode")
			}
			if sceneDomain.Camera.Tilt != 0 {
				camera.Tilt = sceneDomain.Camera.Tilt
				cols = append(cols, "tilt")
			}
		}
		_, err = sess.Cols(cols...).Insert(camera)
		if err != nil {
			logger.Errorln(err)
			return
		}
	}
	{
		cols := []string{"id", "scene_id"}
		fog := &model.SceneFog{
			Id:      util.SnowflakeId(),
			SceneId: scene.Id,
		}
		if sceneDomain.Fog != nil {
			if sceneDomain.Fog.Mode != "" {
				fog.Mode = sceneDomain.Fog.Mode
				cols = append(cols, "mode")
			}
			if sceneDomain.Fog.EndDistance != 1 {
				fog.EndDistance = sceneDomain.Fog.EndDistance
				cols = append(cols, "end_distance")
			}
			if sceneDomain.Fog.StartDistance != 0 {
				fog.StartDistance = sceneDomain.Fog.StartDistance
				cols = append(cols, "start_distance")
			}
			if sceneDomain.Fog.Density != 1 {
				fog.Density = sceneDomain.Fog.Density
				cols = append(cols, "density")
			}
			if sceneDomain.Fog.Enable {
				fog.Enable = 1
				cols = append(cols, "enable")
			}
			if sceneDomain.Fog.Color != nil {
				fog.Color = strings.Join([]string{
					strconv.Itoa(sceneDomain.Fog.Color.Red),
					strconv.Itoa(sceneDomain.Fog.Color.Green),
					strconv.Itoa(sceneDomain.Fog.Color.Blue),
					strconv.Itoa(sceneDomain.Fog.Color.Alpha),
				}, ",")
				cols = append(cols, "color")
			}
		}
		_, err = sess.Cols(cols...).Insert(fog)
		if err != nil {
			logger.Errorln(err)
			return
		}
	}
	err = sess.Commit()
	if err != nil {
		logger.Errorln(err)
		return
	}
	return
}

func (*sceneService) Update(scene *domain.Scene) error {
	if scene.Id == 0 {
		return errors.New("ID must be provided")
	}

	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	var omitCols []string
	{
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
		if !scene.ScaleLegendVisible || scene.Scene.ScaleLegendVisible == 0 {
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
	}

	if _, err := sess.Omit(omitCols...).ID(scene.Id).Update(scene); err != nil {
		logger.Errorln(err)
		return err
	}

	if scene.Atmosphere != nil {
		atmosphere := &model.SceneAtmosphere{
			Visible: 1,
		}
		if !scene.Atmosphere.Visible {
			atmosphere.Visible = -1
		}
		if _, err := sess.Update(atmosphere, &model.SceneAtmosphere{SceneId: scene.Id}); err != nil {
			logger.Errorln(err)
			return err
		}
	}
	if scene.LatLonGrid != nil {
		latLonGrid := &model.SceneLatLonGrid{
			Visible:     -1,
			TextVisible: 1,
		}
		if scene.LatLonGrid.Visible {
			latLonGrid.Visible = 1
		}
		if !scene.LatLonGrid.TextVisible {
			latLonGrid.TextVisible = -1
		}
		if _, err := sess.Update(latLonGrid, &model.SceneLatLonGrid{SceneId: scene.Id}); err != nil {
			logger.Errorln(err)
			return err
		}
	}
	if scene.Camera != nil {
		camera := &model.SceneCamera{
			Altitude:     scene.Camera.Altitude,
			Latitude:     scene.Camera.Latitude,
			Longitude:    scene.Camera.Longitude,
			Heading:      scene.Camera.Heading,
			AltitudeMode: scene.Camera.AltitudeMode,
			Tilt:         scene.Camera.Tilt,
		}
		if _, err := sess.Update(camera, &model.SceneCamera{SceneId: scene.Id}); err != nil {
			logger.Errorln(err)
			return err
		}
	}
	if scene.Fog != nil {
		enable := -1
		if !scene.Fog.Enable {
			enable = 1
		}
		fog := &model.SceneFog{
			Mode:          scene.Fog.Mode,
			EndDistance:   scene.Fog.EndDistance,
			StartDistance: scene.Fog.StartDistance,
			Density:       scene.Fog.Density,
			Enable:        enable,
			Color: strings.Join([]string{
				strconv.Itoa(scene.Fog.Color.Red),
				strconv.Itoa(scene.Fog.Color.Green),
				strconv.Itoa(scene.Fog.Color.Blue),
				strconv.Itoa(scene.Fog.Color.Alpha),
			}, ","),
		}
		if _, err := sess.Update(fog, &model.SceneFog{SceneId: scene.Id}); err != nil {
			logger.Errorln(err)
			return err
		}
	}

	if err := sess.Commit(); err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
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

func (*sceneService) List(condition *model.Scene, offset, limit int, orderBy string) (total int64, list []*model.Scene, err error) {
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
		Scene:              *scene,
		ScaleLegendVisible: scene.ScaleLegendVisible == 1,
		Layers:             nil,
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
