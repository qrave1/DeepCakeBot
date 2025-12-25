package openweather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/qrave1/DeepCakeBot/internal/client/openweather/dto"
)

// OpenWeatherClient клиент для работы с OpenWeather API
type OpenWeatherClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewOpenWeatherClient создает новый клиент для OpenWeather API
func NewOpenWeatherClient(apiKey string) *OpenWeatherClient {
	return &OpenWeatherClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.openweathermap.org/data/2.5",
	}
}

// WeatherData содержит информацию о погоде
type WeatherData struct {
	Temperature float64
	FeelsLike   float64
	Description string
	Humidity    int
	WindSpeed   float64
	Rain        bool
	Snow        bool
}

// GetCurrentWeather получает текущую погоду для указанного города
func (c *OpenWeatherClient) GetCurrentWeather(ctx context.Context, city, countryCode string) (*WeatherData, error) {
	url := fmt.Sprintf(
		"%s/weather?q=%s,%s&appid=%s&units=metric&lang=ru",
		c.baseURL,
		city,
		countryCode,
		c.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("weather API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openWeatherResponse dto.OpenWeatherResponse
	if err := json.Unmarshal(body, &openWeatherResponse); err != nil {
		return nil, fmt.Errorf("failed to parse weather response: %w", err)
	}

	weather := &WeatherData{
		Temperature: openWeatherResponse.Main.Temp,
		FeelsLike:   openWeatherResponse.Main.FeelsLike,
		Humidity:    openWeatherResponse.Main.Humidity,
		WindSpeed:   openWeatherResponse.Wind.Speed,
		Rain:        openWeatherResponse.Rain.OneH > 0,
		Snow:        openWeatherResponse.Snow.OneH > 0,
	}

	if len(openWeatherResponse.Weather) > 0 {
		weather.Description = openWeatherResponse.Weather[0].Description
	}

	return weather, nil
}
