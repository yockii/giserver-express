package domain

type SuperMapMvtIndex struct {
	ViewBounds                SuperMapMvtBound         `json:"viewBounds"`
	Viewer                    SuperMapMvtViewer        `json:"viewer"`
	DistanceUnit              string                   `json:"distanceUnit,omitempty"`
	TileVersion               string                   `json:"tileversion,omitempty"`
	MinVisibleTextSize        int                      `json:"minVisibleTextSize"`
	CoordUnit                 string                   `json:"coordUnit"`
	Scale                     float64                  `json:"scale"`
	Description               string                   `json:"description"`
	PaintBackground           bool                     `json:"paintBackground"`
	MaxVisibleTextSize        int                      `json:"maxVisibleTextSize"`
	MaxVisibleVertex          int                      `json:"maxVisibleVertex"`
	RasterFunction            string                   `json:"rasterfunction"`
	ClipRegionEnabled         bool                     `json:"clipRegionEnabled"`
	TrackingLayer             SuperMapMvtTrackingLayer `json:"trackingLayer"`
	Antialias                 bool                     `json:"antialias"`
	TextOrientationFixed      bool                     `json:"textOrientationFixed"`
	Layers                    []SuperMapMvtLayer       `json:"layers"`
	Angle                     int                      `json:"angle"`
	PrjCoordSys               SuperMapMvtPrjCoorSys    `json:"prjCoordSys"`
	MinScale                  float64                  `json:"minScale"`
	MarkerAngleFixed          bool                     `json:"markerAngleFixed"`
	OverlapDisplayedOptions   struct{}                 `json:"overlapDisplayedOptions"`
	VisibleScales             []float64                `json:"visibleScales"`
	Dpi                       int                      `json:"dpi"`
	VisibleScalesEnabled      bool                     `json:"visibleScalesEnabled"`
	CustomEntireBoundsEnabled bool                     `json:"customEntireBoundsEnabled"`
	ClipRegion                struct{}                 `json:"clipRegion"`
	MaxScale                  float64                  `json:"maxScale"`
	CustomParams              string                   `json:"customParams"`
	Center                    Point                    `json:"center"`
	ColorMode                 struct{}                 `json:"colorMode"`
	TextAngleFixed            bool                     `json:"textAngleFixed"`
	CustomPrjCoordSysType     string                   `json:"customPrjCoordSysType"`
	RectifyType               string                   `json:"rectifyType"`
	OverlapDisplayed          bool                     `json:"overlapDisplayed"`
	UserToken                 SuperMapMvtUserToken     `json:"userToken"`
	CacheEnabled              bool                     `json:"cacheEnabled"`
	DynamicProjection         bool                     `json:"dynamicProjection"`
	AutoAvoidEffectEnabled    bool                     `json:"autoAvoidEffectEnabled"`
	CustomEntireBounds        struct{}                 `json:"customEntireBounds"`
	Name                      string                   `json:"name"`
	Bounds                    SuperMapMvtBound         `json:"bounds"`
	BackgroundStyle           struct{}                 `json:"backgroundStyle"`
	ReturnImage               bool                     `json:"returnImage"`
	ReturnType                string                   `json:"returnType"`
}

type SuperMapMvtUserToken struct {
	UserID string `json:"userID"`
}

type SuperMapMvtPrjCoorSys struct {
	DistanceUnit    string                 `json:"distanceUnit"`
	ProjectionParam SuperMapMvtProjParam   `json:"projectionParam"`
	EpsgCode        string                 `json:"epsgCode"`
	CoordUnit       string                 `json:"coordUnit"`
	Name            string                 `json:"name"`
	Projection      SuperMapMvtProjection  `json:"projection"`
	Type            string                 `json:"type"`
	CoordSystem     SuperMapMvtCoordSystem `json:"coordSystem"`
}

type SuperMapMvtCoordSystem struct {
	Datum          SuperMapMvtCoordSystemDatum         `json:"datum"`
	Unit           string                              `json:"unit"`
	SpatialRefType string                              `json:"spatialRefType"`
	Name           string                              `json:"name"`
	Type           string                              `json:"type"`
	PrimeMeridian  SuperMapMvtCoordSystemPrimeMeridian `json:"primeMeridian"`
}

type SuperMapMvtCoordSystemPrimeMeridian struct {
	LongitudeValue int    `json:"longitudeValue"`
	Name           string `json:"name"`
	Type           string `json:"type"`
}

type SuperMapMvtCoordSystemDatum struct {
	Name     string                              `json:"name"`
	Type     string                              `json:"type"`
	Spheroid SuperMapMvtCoordSystemDatumSpheroid `json:"spheroid"`
}

type SuperMapMvtCoordSystemDatumSpheroid struct {
	Flatten float64 `json:"flatten"`
	Name    string  `json:"name"`
	Axis    float64 `json:"axis"`
	Type    string  `json:"type"`
}

type SuperMapMvtProjection struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SuperMapMvtProjParam struct {
	CentralParallel        int     `json:"centralParallel"`
	FirstPointLongitude    float64 `json:"firstPointLongitude"`
	RectifiedAngle         int     `json:"rectifiedAngle"`
	ScaleFactor            int     `json:"scaleFactor"`
	FalseNorthing          int     `json:"falseNorthing"`
	CentralMeridian        int     `json:"centralMeridian"`
	SecondStandardParallel int     `json:"secondStandardParallel"`
	SecondPointLongitude   int     `json:"secondPointLongitude"`
	Azimuth                int     `json:"azimuth"`
	FalseEasting           int     `json:"falseEasting"`
	FirstStandardParallel  int     `json:"firstStandardParallel"`
}

type SuperMapMvtLayer struct {
	QueryAble   bool              `json:"queryable"`
	Visible     bool              `json:"visible"`
	Name        string            `json:"name"`
	Bounds      SuperMapMvtBound  `json:"bounds"`
	Caption     string            `json:"caption"`
	Description string            `json:"description"`
	SubLayers   *SuperMapMvtLayer `json:"subLayers"`
	Type        string            `json:"type"`
}

type SuperMapMvtTrackingLayer struct {
	HighlightTargets []string `json:"highlightTargets"`
}

type SuperMapMvtBound struct {
	Top        float64 `json:"top"`
	Left       float64 `json:"left"`
	Bottom     float64 `json:"bottom"`
	Right      float64 `json:"right"`
	LeftBottom Point   `json:"leftBottom"`
	RightTop   Point   `json:"rightTop"`
}

type SuperMapMvtViewer struct {
	Top         float64 `json:"top"`
	Left        float64 `json:"left"`
	Bottom      float64 `json:"bottom"`
	Right       float64 `json:"right"`
	Width       float64 `json:"width"`
	Height      float64 `json:"height"`
	LeftTop     Point   `json:"leftTop"`
	RightBottom Point   `json:"rightBottom"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
