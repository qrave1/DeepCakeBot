package dto

// OpenWeatherResponse структура ответа от OpenWeather API
type OpenWeatherResponse struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
		Main        string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain,omitempty"`
	Snow struct {
		OneH float64 `json:"1h"`
	} `json:"snow,omitempty"`
}
