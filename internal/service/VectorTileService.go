package service

import (
	"encoding/json"
	"errors"
	"github.com/beevik/etree"
	logger "github.com/sirupsen/logrus"
	"github.com/yockii/giserver-express/internal/domain"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/util"
)

var VectorTileService = &vectorTileService{
	indexCache: make(map[string]string),
	vtInfoMap:  make(map[string]*model.VectorTile),
}

type vectorTileService struct {
	indexCache map[string]string
	vtInfoMap  map[string]*model.VectorTile
	lock       sync.Mutex
}

func (s *vectorTileService) Add(vt *model.VectorTile) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.VectorTile{
		Name: vt.Name,
	})
	if err != nil {
		logger.Errorln(err)
		return
	}
	if c > 0 {
		duplicated = true
		return
	}

	vt.Id = database.Int64(util.SnowflakeId())

	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()
	_, err = sess.Insert(vt)
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

func (*vectorTileService) Update(vectorTile *model.VectorTile) error {
	if vectorTile.Id == 0 {
		return errors.New("ID must be provided")
	}

	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	if _, err := sess.ID(vectorTile.Id).Update(vectorTile); err != nil {
		logger.Errorln(err)
		return err
	}

	if err := sess.Commit(); err != nil {
		logger.Errorln(err)
		return err
	}
	return nil
}

func (*vectorTileService) Delete(id database.Int64) (err error) {
	sess := database.DB.NewSession()
	sess.Begin()
	defer sess.Close()

	_, err = sess.Delete(&model.VectorTile{Id: id})
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

func (*vectorTileService) List(condition *model.VectorTile, offset, limit int, orderBy string) (total int64, list []*model.VectorTile, err error) {
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

func (s *vectorTileService) FindByName(name string) (*model.VectorTile, error) {
	if name == "" {
		return nil, nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if vt, ok := s.vtInfoMap[name]; ok && vt != nil {
		return vt, nil
	}

	vectorTile := new(model.VectorTile)
	vectorTile.Name = name
	exist, err := database.DB.Get(vectorTile)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	s.vtInfoMap[name] = vectorTile

	return vectorTile, nil
}

func (s *vectorTileService) GetIndex(name string) (string, error) {
	if name == "" {
		return "", nil
	}
	if ij, ok := s.indexCache[name]; ok && ij != "" {
		return ij, nil
	}

	vectorTile, err := s.FindByName(name)
	if err != nil {
		return "", err
	}
	if vectorTile == nil {
		return "", nil
	}

	// 存在则获取对应的文件
	store, err := StoreService.GetById(vectorTile.StoreId)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}
	if store == nil {
		return "", nil
	}
	xmlDoc := etree.NewDocument()
	if store.StoreType == -1 {
		// 本地存储
		err = xmlDoc.ReadFromFile(filepath.Join(store.Path, vectorTile.PathName, vectorTile.IndexFileName))
		if err != nil {
			logger.Errorln(err)
			return "", err
		}
	} else if store.StoreType == 1 {
		objectKey := store.Path + vectorTile.PathName + "/" + vectorTile.IndexFileName
		var reader io.Reader
		reader, err = OssService.StreamFromStore(store, objectKey)
		_, err = xmlDoc.ReadFrom(reader)
		if err != nil {
			return "", err
		}
	} else {
		return "", nil
	}
	// xml读取完毕，开始构造json
	root := xmlDoc.SelectElement("SuperMapCache")

	scaleValueElements := root.FindElements("//sml:Scales/sml:Scale/sml:Value")
	var scales []float64
	for _, v := range scaleValueElements {
		t := v.Text()
		f, e := strconv.ParseFloat(t, 64)
		if e != nil {
			logger.Errorln(e)
			return "", e
		}
		scales = append(scales, f)
	}

	smlBounds := root.SelectElement("sml:Bounds")
	bound := domain.SuperMapMvtBound{}
	{
		smlBoundsLeft := smlBounds.SelectElement("sml:Left")
		bound.Left, err = strconv.ParseFloat(smlBoundsLeft.Text(), 64)
		if err != nil {
			logger.Errorln(err)
			return "", err
		}
		smlBoundsTop := smlBounds.SelectElement("sml:Top")
		bound.Top, err = strconv.ParseFloat(smlBoundsTop.Text(), 64)
		if err != nil {
			logger.Errorln(err)
			return "", err
		}
		smlBoundsRight := smlBounds.SelectElement("sml:Right")
		bound.Right, err = strconv.ParseFloat(smlBoundsRight.Text(), 64)
		if err != nil {
			logger.Errorln(err)
			return "", err
		}
		smlBoundsBottom := smlBounds.SelectElement("sml:Bottom")
		bound.Bottom, err = strconv.ParseFloat(smlBoundsBottom.Text(), 64)
		if err != nil {
			logger.Errorln(err)
			return "", err
		}
		bound.LeftBottom = domain.Point{
			X: bound.Left,
			Y: bound.Bottom,
		}
		bound.RightTop = domain.Point{
			X: bound.Right,
			Y: bound.Top,
		}
	}

	smlImageSizeElement := root.SelectElement("sml:ImageSize")
	var imageSize float64 = 512
	imageSize, err = strconv.ParseFloat(smlImageSizeElement.Text(), 64)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}

	idx := &domain.SuperMapMvtIndex{
		ViewBounds: bound,
		Viewer: domain.SuperMapMvtViewer{
			Top:     0,
			Left:    0,
			Bottom:  imageSize,
			Right:   imageSize,
			Width:   imageSize,
			Height:  imageSize,
			LeftTop: domain.Point{},
			RightBottom: domain.Point{
				X: imageSize,
				Y: imageSize,
			},
		},
		CoordUnit: "DEGREE",
		Scale:     0,
		TrackingLayer: domain.SuperMapMvtTrackingLayer{
			HighlightTargets: make([]string, 0),
		},
		VisibleScales:        scales,
		Dpi:                  96,
		VisibleScalesEnabled: true,
		Center: domain.Point{
			X: (bound.Left + bound.Right) / 2,
			Y: (bound.Top + bound.Bottom) / 2,
		},
		RectifyType:            "BYCENTERANDMAPSCALE",
		UserToken:              domain.SuperMapMvtUserToken{},
		AutoAvoidEffectEnabled: true,
		Name:                   vectorTile.Name,
		Bounds:                 bound,
		ReturnType:             "URL",
	}
	{
		scaleOriginalElement := root.SelectElement("sml:ScaleOriginalResolution")
		idx.Scale, err = strconv.ParseFloat(scaleOriginalElement.Text(), 64)

		idx.Layers = append(idx.Layers, domain.SuperMapMvtLayer{
			Visible: true,
			Name:    vectorTile.Name,
			Bounds:  bound,
			Type:    "CUSTOM",
		})

		coordinateReferenceSystemElement := root.SelectElement("sml:CoordinateReferenceSystem")
		distanceUnitElement := coordinateReferenceSystemElement.SelectElement("sml:Units")
		ePSGCodeElement := coordinateReferenceSystemElement.SelectElement("sml:EPSGCode")
		geographicCoordinateSystemElement := coordinateReferenceSystemElement.SelectElement("sml:GeographicCoordinateSystem")
		CoordUnitElement := geographicCoordinateSystemElement.SelectElement("sml:Units")
		NamesetElement := coordinateReferenceSystemElement.SelectElement("sml:Nameset")
		nameElement := NamesetElement.SelectElement("sml:name")

		horizonalGeodeticDatumElement := geographicCoordinateSystemElement.SelectElement("sml:HorizonalGeodeticDatum")
		datumNamesetElement := horizonalGeodeticDatumElement.SelectElement("sml:Nameset")
		datumNameElement := datumNamesetElement.SelectElement("sml:Name")
		ellipsoidElement := horizonalGeodeticDatumElement.SelectElement("sml:Ellipsoid")
		spheroidNamesetElement := ellipsoidElement.SelectElement("sml:Nameset")
		spheroidNameElement := spheroidNamesetElement.SelectElement("sml:Name")
		semiMajorAxisElement := ellipsoidElement.SelectElement("sml:SemiMajorAxis")
		var semiMajorAxis float64
		semiMajorAxis, err = strconv.ParseFloat(semiMajorAxisElement.Text(), 64)
		if err != nil {
			logger.Errorln(err)
			return "", err
		}

		primeMeridianElement := geographicCoordinateSystemElement.SelectElement("sml:PrimeMeridian")
		primeMeridianNamesetElement := primeMeridianElement.SelectElement("sml:Nameset")
		primeMeridianNameElement := primeMeridianNamesetElement.SelectElement("sml:Name")

		idx.PrjCoordSys = domain.SuperMapMvtPrjCoorSys{
			DistanceUnit:    distanceUnitElement.Text(),
			ProjectionParam: domain.SuperMapMvtProjParam{},
			EpsgCode:        ePSGCodeElement.Text(),
			CoordUnit:       CoordUnitElement.Text(),
			Name:            nameElement.Text(),
			Projection: domain.SuperMapMvtProjection{
				Type: "PRJ_NONPROJECTION",
			},
			Type: "PCS_EARTH_LONGITUDE_LATITUDE",
			CoordSystem: domain.SuperMapMvtCoordSystem{
				Datum: domain.SuperMapMvtCoordSystemDatum{
					Name: datumNameElement.Text(),
					Type: "DATUM_CHINA_2000",
					Spheroid: domain.SuperMapMvtCoordSystemDatumSpheroid{
						Flatten: 0.003352810681182319,
						Name:    spheroidNameElement.Text(),
						Axis:    semiMajorAxis,
						Type:    "SPHEROID_CHINA_2000",
					},
				},
				Unit:           CoordUnitElement.Text(),
				SpatialRefType: "SPATIALREF_EARTH_PROJECTION",
				Name:           nameElement.Text(),
				Type:           strings.ToUpper(nameElement.Text()),
				PrimeMeridian: domain.SuperMapMvtCoordSystemPrimeMeridian{
					LongitudeValue: 0,
					Name:           primeMeridianNameElement.Text(),
					Type:           "PRIMEMERIDIAN_" + strings.ToUpper(primeMeridianNameElement.Text()),
				},
			},
		}
	}
	var result []byte
	result, err = json.Marshal(idx)
	if err != nil {
		logger.Errorln(err)
		return "", err
	}

	indexStr := string(result)
	s.lock.Lock()
	defer s.lock.Unlock()
	s.indexCache[name] = indexStr
	return indexStr, nil
}

func (s *vectorTileService) ClearCache(names ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, name := range names {
		delete(s.vtInfoMap, name)
		delete(s.indexCache, name)
	}
}

func (s *vectorTileService) ReadMvtFile(name string, dirAndFile ...string) (io.Reader, error) {
	vt, err := s.FindByName(name)
	if err != nil {
		return nil, err
	}
	if vt == nil {
		return nil, nil
	}
	store, err := StoreService.GetById(vt.StoreId)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if store == nil {
		return nil, nil
	}
	if store.StoreType == -1 {
		// 本地存储
		var file *os.File
		fp := []string{
			store.Path,
			vt.PathName,
		}
		fp = append(fp, dirAndFile...)
		file, err = os.Open(filepath.Join(fp...))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil, nil
			}
			if e, ok := err.(*fs.PathError); ok && e.Err == syscall.ENOSPC {
				return nil, nil
			} else {
				logger.Errorln(err)
				return nil, err
			}
		}
		return file, nil
	} else if store.StoreType == 1 {
		objectKey := store.Path + vt.PathName
		for _, p := range dirAndFile {
			objectKey += "/" + p
		}
		var reader io.Reader
		reader, err = OssService.StreamFromStore(store, objectKey)
		if err != nil {
			return nil, err
		}
		return reader, nil
	}
	return nil, nil
}
