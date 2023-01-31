package model

import "github.com/yockii/giserver-express/pkg/database"

type Space struct {
	Id         int64             `json:"id" xorm:"pk"`
	Name       string            `json:"name" xorm:"index"` //
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func init() {
	SyncModels = append(SyncModels, Space{})
}
