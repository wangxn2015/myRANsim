package model

// MapLayout represents information required for geo-map visualizations
type MapLayout struct {
	Center         Coordinate `mapstructure:"center"`
	Zoom           float32    `mapstructure:"zoom"`
	LocationsScale float32    `mapstructure:"locationsScale"`
	FadeMap        bool       `mapstructure:"fade"`
	ShowRoutes     bool       `mapstructure:"showRoutes"`
	ShowPower      bool       `mapstructure:"showPower"`
}
