package service

import (
	"errors"
	logger "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/yockii/giserver-express/internal/model"
	"github.com/yockii/giserver-express/pkg/database"
	"github.com/yockii/giserver-express/pkg/util"
)

var MapTileService = &mapTileService{
	mtInfoMap: make(map[string]*model.MapTile),
}

type mapTileService struct {
	mtInfoMap map[string]*model.MapTile
	lock      sync.Mutex
}

func (s *mapTileService) Add(vt *model.MapTile) (duplicated bool, err error) {
	var c int64
	c, err = database.DB.Count(&model.MapTile{
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

func (*mapTileService) Update(vectorTile *model.MapTile) error {
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

func (*mapTileService) Delete(id database.Int64) (err error) {
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

func (*mapTileService) List(condition *model.MapTile, offset, limit int, orderBy string) (total int64, list []*model.MapTile, err error) {
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

func (s *mapTileService) FindByName(name string) (*model.MapTile, error) {
	if name == "" {
		return nil, nil
	}
	s.lock.Lock()
	defer s.lock.Unlock()

	if vt, ok := s.mtInfoMap[name]; ok && vt != nil {
		return vt, nil
	}

	vectorTile := new(model.MapTile)
	vectorTile.Name = name
	exist, err := database.DB.Get(vectorTile)
	if err != nil {
		logger.Errorln(err)
		return nil, err
	}
	if !exist {
		return nil, nil
	}

	s.mtInfoMap[name] = vectorTile

	return vectorTile, nil
}

func (s *mapTileService) ClearCache(names ...string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, name := range names {
		delete(s.mtInfoMap, name)
	}
}

func (s *mapTileService) ReadFile(name string, dirAndFile ...string) (io.Reader, error) {
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
			var e *os.PathError
			if errors.As(err, &e) && e.Err == syscall.ENOSPC {
				return nil, nil
			}
		}
		return file, nil
	} else if store.StoreType == 1 {
		objectKey := store.Path + vt.PathName
		for _, p := range dirAndFile {
			if len(p) > 0 {
				objectKey += "/" + p
			}
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
