package model

import "github.com/yockii/giserver-express/pkg/database"

// Data 在Space之下
type Data struct {
	Id              int64             `json:"id,omitempty" xorm:"pk"`
	SpaceId         int64             `json:"spaceId,omitempty" xorm:"index"`
	Name            string            `json:"name" xorm:"default('Config')"`
	DataType        string            `json:"dataType" xorm:"default('OSGB')"` // S3M
	DataConfigPath  string            `json:"dataConfigPath" xorm:"comment('scp存放位置')"`
	DataStoreTypeId int64             `json:"dataStoreTypeId" xorm:"comment('数据存放类型id 0-本地文件')"`
	CreateTime      database.DateTime `json:"createTime" xorm:"created"`
}

func init() {
	SyncModels = append(SyncModels, Data{})
}
