package usecase

import (
	"context"
	"fmt"

	"github.com/qrave1/DeepCakeBot/internal/storage"

	tele "gopkg.in/telebot.v3"
)

type ApplicationBot struct {
	bot            *tele.Bot
	storage        storage.UserRepository
	weatherService *WeatherService
}

// NewApplicationBot создает новый сервис бота
func NewApplicationBot(bot *tele.Bot, storage storage.UserRepository, weatherService *WeatherService) *ApplicationBot {
	return &ApplicationBot{
		bot:            bot,
		storage:        storage,
		weatherService: weatherService,
	}
}

// RegisterHandlers регистрирует все обработчики команд и callback'ов
func (s *ApplicationBot) RegisterHandlers() {
	s.bot.Handle("/weather", s.handleGetWeather)

	// Обработчик команды /start
	s.bot.Handle("/start", s.handleStart)

	// Обработчик команды /settings
	s.bot.Handle("/settings", s.handleSettings)

	// Обработчики callback для настроек
	s.bot.Handle(&btnEnableWeather, s.handleEnableWeather)
	s.bot.Handle(&btnDisableWeather, s.handleDisableWeather)
}

// SendWeatherToUser отправляет прогноз погоды конкретному пользователю
func (s *ApplicationBot) SendWeatherToUser(ctx context.Context, chatID int64) error {
	weather, err := s.weatherService.GetWeather(ctx)
	if err != nil {
		return fmt.Errorf("failed to get weather: %w", err)
	}

	message := s.weatherService.FormatWeatherMessage(weather)

	_, err = s.bot.Send(&tele.Chat{ID: chatID}, message)
	if err != nil {
		return fmt.Errorf("failed to send message to %d: %w", chatID, err)
	}

	return nil
}
