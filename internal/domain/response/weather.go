package response

type MeteosourceWind struct {
	Speed float32 `json:"speed"`
	Angle int     `json:"angle"`
	Dir   string  `json:"dir"`
}

type MeteosourcePrecipitation struct {
	Total float32 `json:"total"`
	Type  string  `json:"type"`
}

type MeteosourceWeatherData struct {
	Icon          string                   `json:"icon"`
	IconNum       int                      `json:"icon_num"`
	Summary       string                   `json:"summary"`
	Temperature   float32                  `json:"temperature"`
	CloudCover    int                      `json:"cloud_cover"`
	Wind          MeteosourceWind          `json:"wind"`
	Precipitation MeteosourcePrecipitation `json:"precipitation"`
}

type MeteosourceWeather struct {
	Lat       string                 `json:"lat"`
	Lon       string                 `json:"lon"`
	Elevation int                    `json:"elevation"`
	Timezone  string                 `json:"timezone"`
	Units     string                 `json:"units"`
	Current   MeteosourceWeatherData `json:"current"`
}
