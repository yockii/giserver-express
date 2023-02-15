package model

import "github.com/yockii/giserver-express/pkg/database"

// Data 在Space之下
type Data struct {
	Id              int64             `json:"id,omitempty" xorm:"pk"`
	SpaceId         int64             `json:"spaceId,omitempty" xorm:"index"`
	Name            string            `json:"name" xorm:"default('Config')"`
	DataType        string            `json:"dataType" xorm:"default('OSGB')"` // S3M/KML等
	DataConfigPath  string            `json:"dataConfigPath" xorm:"comment('osgb/s3m等格式为scp存放位置，KML等则为自身文件存放位置')"`
	DataStoreTypeId int64             `json:"dataStoreTypeId" xorm:"comment('数据存放类型id 0-本地文件')"`
	CreateTime      database.DateTime `json:"createTime" xorm:"created"`
}

func init() {
	SyncModels = append(SyncModels, Data{})
}
