package model

import "github.com/yockii/giserver-express/pkg/database"

type Store struct {
	Id              database.Int64    `json:"id,omitempty" xorm:"pk"`
	Name            string            `json:"name,omitempty"`
	StoreType       int               `json:"storeType" xorm:"default(1) comment('存储类型 0-本地文件 1-oss')"`
	Endpoint        string            `json:"endpoint,omitempty"`
	AccessKeyId     string            `json:"accessKeyId,omitempty"`
	AccessKeySecret string            `json:"accessKeySecret,omitempty"`
	BucketName      string            `json:"bucketName,omitempty"`
	Path            string            `json:"path"`
	UseCname        int               `json:"useCname" xorm:"default(0) comment('使用自定义域名，0-否，1-是')"`
	CreateTime      database.DateTime `json:"createTime" xorm:"created"`
}

func init() {
	SyncModels = append(SyncModels, Store{})
}
