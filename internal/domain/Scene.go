package domain

import "github.com/yockii/giserver-express/internal/model"

type SpaceScene struct {
	ResourceConfigId    string   `json:"resourceConfigID,omitempty"`
	SupportedMediaTypes []string `json:"supportedMediaTypes,omitempty"`
	Path                string   `json:"path,omitempty"`
	Name                string   `json:"name,omitempty"`
	ResourceType        string   `json:"resourceType,omitempty"`
}

type Scene struct {
	model.Scene
	MinCameraDistance   float64          `json:"minCameraDistance" `
	MaxCameraDistance   float64          `json:"maxCameraDistance" `
	ScaleLegendVisible  bool             `json:"scaleLegendVisible"`
	CameraFOV           float64          `json:"cameraFOV" `
	FogVisibleAltitude  int64            `json:"fogVisibleAltitude"`
	SceneType           string           `json:"sceneType"`
	TerrainExaggeration int              `json:"terrainExaggeration"`
	TrackingLayer       *string          `json:"trackingLayer"`
	Xml                 *string          `json:"xml"`
	ScreenLayer         *string          `json:"screenLayer"`
	Atmosphere          *SceneAtmosphere `json:"atmosphere"`
	LatLonGrid          *SceneLatLonGrid `json:"latLonGrid"`
	Layers              []*SceneLayer    `json:"layers"`
	Camera              *SceneCamera     `json:"camera"`
	Fog                 *SceneFog        `json:"fog"`
}

type SceneAtmosphere struct {
	Visible bool `json:"visible"`
}

type SceneLatLonGrid struct {
	Visible     bool `json:"visible" `
	TextVisible bool `json:"textVisible"`
}

type SceneLayer struct {
	Name               string      `json:"name"`
	MinVisibleAltitude float64     `json:"minVisibleAltitude"`
	MaxVisibleAltitude float64     `json:"maxVisibleAltitude"`
	VisibleDistance    float64     `json:"visibleDistance"`
	IsWebDatasource    bool        `json:"isWebDatasource"`
	AlwaysRender       bool        `json:"alwaysRender"`
	Visible            bool        `json:"visible"`
	Level              int         `json:"level"`
	UseTwoDimenCache   bool        `json:"useTwoDimenCache"`
	Editable           bool        `json:"editable"`
	Caption            string      `json:"caption"`
	Description        string      `json:"description"`
	DataName           string      `json:"dataName"`
	HasLocalCache      bool        `json:"hasLocalCache,omitempty"`
	Layer3DType        string      `json:"layer3DType"`
	DataConfigPath     string      `json:"dataConfigPath"`
	ExtendXML          interface{} `json:"extendXML"`
	SubLayers          interface{} `json:"subLayers"`
	Type               interface{} `json:"type"`
	ParentLayerName    interface{} `json:"parentLayerName"`
	Bounds             interface{} `json:"bounds"`
	CachePassword      interface{} `json:"cachePassword"`
	Queryable          bool        `json:"queryable"`
	OldCache           bool        `json:"oldCache"`
}

type SceneCamera struct {
	Altitude     float64 `json:"altitude" `
	Latitude     float64 `json:"latitude" `
	Longitude    float64 `json:"longitude" `
	Heading      float64 `json:"heading" `
	AltitudeMode string  `json:"altitudeMode" `
	Tilt         float64 `json:"tilt"`
	Empty        int     `json:"empty"`
}

type SceneFog struct {
	Mode          string  `json:"mode"`
	EndDistance   float64 `json:"endDistance" `
	StartDistance float64 `json:"startDistance"`
	Color         *Color  `json:"color"`
	Density       float64 `json:"density"`
	Enable        bool    `json:"enable"`
}

type Color struct {
	Red   int `json:"red"`
	Green int `json:"green"`
	Blue  int `json:"blue"`
	Alpha int `json:"alpha"`
}
