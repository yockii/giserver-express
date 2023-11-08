package model

import "github.com/yockii/giserver-express/pkg/database"

type MapTile struct {
	Id         database.Int64    `json:"id" xorm:"pk"`
	Name       string            `json:"name" xorm:"index"`
	StoreId    database.Int64    `json:"storeId"`
	PathName   string            `json:"pathName" xorm:"comment('存储信息的相对路径')"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func init() {
	SyncModels = append(SyncModels, MapTile{})
}
