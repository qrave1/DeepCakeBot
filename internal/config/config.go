package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v10"
)

// Config содержит все настройки приложения
type Config struct {
	// Telegram Bot Token
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`

	// OpenWeather API Key
	OpenWeatherAPIKey string `env:"OPENWEATHER_API_KEY,required"`

	// Database URL для подключения к PostgreSQL
	DatabaseURL string `env:"DATABASE_URL,required"`

	// Timezone для планировщика (по умолчанию Europe/Moscow)
	Timezone string `env:"TIMEZONE" envDefault:"Europe/Moscow"`

	// Время отправки прогноза погоды (по умолчанию 07:00)
	WeatherScheduleHour int `env:"WEATHER_SCHEDULE_HOUR" envDefault:"7"`

	// Город для прогноза погоды
	City string `env:"CITY" envDefault:"Moscow"`

	// Код страны для OpenWeather API
	CountryCode string `env:"COUNTRY_CODE" envDefault:"RU"`
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if _, err := time.LoadLocation(cfg.Timezone); err != nil {
		return nil, fmt.Errorf("invalid timezone %s: %w", cfg.Timezone, err)
	}

	if cfg.WeatherScheduleHour < 0 || cfg.WeatherScheduleHour > 23 {
		return nil, fmt.Errorf("invalid weather schedule hour: %d (must be 0-23)", cfg.WeatherScheduleHour)
	}

	return cfg, nil
}

