package usecase

import (
	"context"
	"fmt"

	"github.com/qrave1/DeepCakeBot/internal/client/openweather"
)

// WeatherService предоставляет информацию о погоде и рекомендации
type WeatherService struct {
	client      *openweather.OpenWeatherClient
	city        string
	countryCode string
}

// NewWeatherService создает новый сервис погоды
func NewWeatherService(apiKey, city, countryCode string) *WeatherService {
	return &WeatherService{
		client:      openweather.NewOpenWeatherClient(apiKey),
		city:        city,
		countryCode: countryCode,
	}
}

// GetWeather получает текущую погоду для заданного города
func (s *WeatherService) GetWeather(ctx context.Context) (*openweather.WeatherData, error) {
	return s.client.GetCurrentWeather(ctx, s.city, s.countryCode)
}

// GetClothingRecommendation возвращает рекомендации по одежде на основе погоды
func (s *WeatherService) GetClothingRecommendation(weather *openweather.WeatherData) string {
	temp := weather.Temperature
	var recommendation string

	switch {
	case temp < -15:
		recommendation = "🧥 Очень холодно! Теплая зимняя одежда, шапка, шарф, перчатки обязательны."
	case temp >= -15 && temp < -5:
		recommendation = "❄️ Холодно. Зимняя куртка, теплые аксессуары (шапка, перчатки)."
	case temp >= -5 && temp < 5:
		recommendation = "🧥 Прохладно. Демисезонная куртка, можно добавить шарф."
	case temp >= 5 && temp < 15:
		recommendation = "🧥 Прохладная погода. Легкая куртка или толстовка."
	case temp >= 15 && temp < 25:
		recommendation = "👕 Комфортная температура. Легкая одежда, можно без куртки."
	default:
		recommendation = "☀️ Жарко! Легкая летняя одежда, не забудьте солнцезащитные средства."
	}

	if weather.Rain {
		recommendation += "\n☔ Ожидается дождь - возьмите зонт или дождевик!"
	}
	if weather.Snow {
		recommendation += "\n❄️ Ожидается снег - одевайтесь теплее и будьте осторожны на дорогах!"
	}

	return recommendation
}

// FormatWeatherMessage форматирует сообщение с прогнозом погоды
func (s *WeatherService) FormatWeatherMessage(weather *openweather.WeatherData) string {
	msg := fmt.Sprintf(
		"🌤 Прогноз погоды для %s:\n\n"+
			"🌡 Температура: %.1f°C (ощущается как %.1f°C)\n"+
			"📝 Описание: %s\n"+
			"💧 Влажность: %d%%\n"+
			"💨 Скорость ветра: %.1f м/с\n\n"+
			"%s",
		s.city,
		weather.Temperature,
		weather.FeelsLike,
		weather.Description,
		weather.Humidity,
		weather.WindSpeed,
		s.GetClothingRecommendation(weather),
	)

	return msg
}
