package model

import "github.com/yockii/giserver-express/pkg/database"

type Scene struct {
	Id                  database.Int64    `json:"id" xorm:"pk"`
	SpaceId             database.Int64    `json:"spaceId" xorm:"index"`
	Name                string            `json:"name"`
	ResourceConfigId    string            `json:"resourceConfigID" xorm:"default('scene')"`               // scene
	SupportedMediaTypes string            `json:"supportedMediaTypes" xorm:"default('application/json')"` // ["application/xml","text/xml","application/json","application/fastjson","application/rjson","text/html","application/jsonp","application/x-java-serialized-object","application/realspace","application/openrealspace","application/scenezip"]
	ResourceType        string            `json:"resourceType" xorm:"default('ArithmeticResource')"`      // ArithmeticResource
	MinCameraDistance   float64           `json:"minCameraDistance" xorm:"default(6367103)"`
	MaxCameraDistance   float64           `json:"maxCameraDistance" xorm:"default(47836027.5)"`
	ScaleLegendVisible  int               `json:"scaleLegendVisible" xorm:"default(1)"`
	CameraFov           float64           `json:"cameraFOV" xorm:"default(45.0000000001462)"`
	FogVisibleAltitude  int64             `json:"fogVisibleAltitude" xorm:"default(20000)"`
	SceneType           string            `json:"sceneType" xorm:"default('GLOBE')"`
	TerrainExaggeration int               `json:"terrainExaggeration" xorm:"default(1)"`
	CreateTime          database.DateTime `json:"createTime" xorm:"created"`
}

type SceneAtmosphere struct {
	Id      database.Int64 `json:"id,omitempty" xorm:"pk"`
	SceneId database.Int64 `json:"sceneId,omitempty" xorm:"index"`
	Visible int            `json:"visible,omitempty" xorm:"default(1)"`
}

type SceneLatLonGrid struct {
	Id          database.Int64 `json:"id,omitempty" xorm:"pk"`
	SceneId     database.Int64 `json:"sceneId,omitempty" xorm:"index"`
	Visible     int            `json:"visible,omitempty" xorm:"default(0)"`
	TextVisible int            `json:"textVisible,omitempty" xorm:"default(1)"`
}

type SceneCamera struct {
	Id           database.Int64 `json:"id,omitempty" xorm:"pk"`
	SceneId      database.Int64 `json:"sceneId,omitempty" xorm:"index"`
	Altitude     float64        `json:"altitude" xorm:"default(312.24457270652056)"`
	Latitude     float64        `json:"latitude" xorm:"default(45.766272055085)"`
	Longitude    float64        `json:"longitude" xorm:"default(126.62133125438919)"`
	Heading      float64        `json:"heading" xorm:"default(0.6557127322272999)"`
	AltitudeMode string         `json:"altitudeMode" xorm:"default('ABSOLUTE')"`
	Tilt         float64        `json:"tilt" xorm:"default(51.599254528603325)"`
	Empty        int            `json:"empty" xorm:"default(-1)"`
}

type SceneFog struct {
	Id            database.Int64 `json:"id,omitempty" xorm:"pk"`
	SceneId       database.Int64 `json:"sceneId,omitempty" xorm:"index"`
	Mode          string         `json:"mode" xorm:"default('EXP')"`
	EndDistance   float64        `json:"end_distance" xorm:"default(1)"`
	StartDistance float64        `json:"start_distance" xorm:"default(0)"`
	Color         string         `json:"color" xorm:"default('255,255,255,255')"`
	Density       float64        `json:"density" xorm:"default(1)"`
	Enable        int            `json:"enable" xorm:"default(-1)"`
}

type SceneLayer struct {
	Id                 database.Int64    `json:"id,omitempty" xorm:"pk"`
	SpaceId            database.Int64    `json:"spaceId,omitempty"`
	SceneId            database.Int64    `json:"sceneId,omitempty" xorm:"index"`
	DataId             database.Int64    `json:"dataId"`
	Name               string            `json:"name" xorm:"default('Config')"`
	MinVisibleAltitude float64           `json:"minVisibleAltitude" xorm:"default(0)"`
	MaxVisibleAltitude float64           `json:"maxVisibleAltitude" xorm:"default(0)"`
	VisibleDistance    float64           `json:"visibleDistance" xorm:"default(0)"`
	IsWebDatasource    int               `json:"isWebDatasource" xorm:"default(0)"`
	AlwaysRender       int               `json:"alwaysRender" xorm:"default(1)"`
	Visible            int               `json:"visible,omitempty" xorm:"default(0)"`
	Level              int               `json:"level,omitempty" xorm:"default(-1)"`
	UseTwoDimenCache   int               `json:"useTwoDimenCache,omitempty" xorm:"default(0)"`
	Editable           int               `json:"editable,omitempty" xorm:"default(0)"`
	Caption            string            `json:"caption" xorm:"default('Config')"`
	Description        string            `json:"description"`
	HasLocalCache      int               `json:"hasLocalCache,omitempty" xorm:"default(1)"`
	Layer3DType        string            `json:"layer3DType" xorm:"layer_3d_type default('OSGBLayer')"`
	CreateTime         database.DateTime `json:"createTime" xorm:"created"`
}

func init() {
	SyncModels = append(SyncModels,
		Scene{},
		SceneAtmosphere{},
		SceneLatLonGrid{},
		SceneCamera{},
		SceneFog{},
		SceneLayer{},
	)
}
