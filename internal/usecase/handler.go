package usecase

import (
	"context"
	"log"

	tele "gopkg.in/telebot.v3"
)

// handleStart обрабатывает команду /start
func (s *ApplicationBot) handleStart(c tele.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	// Создаем или обновляем пользователя
	if err := s.storage.CreateUser(ctx, chatID); err != nil {
		log.Printf("Failed to create user %d: %v", chatID, err)
		return c.Send("Произошла ошибка при регистрации. Попробуйте позже.")
	}

	welcomeMsg := "👋 Добро пожаловать в DeepCake Bot!\n\n" +
		"Я буду отправлять вам прогноз погоды каждое утро в 07:00 по МСК.\n\n" +
		"Доступные команды:\n" +
		"/settings - настройки рассылки"

	return c.Send(welcomeMsg)
}

func (s *ApplicationBot) handleGetWeather(c tele.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	err := s.SendWeatherToUser(ctx, chatID)
	if err != nil {
		return err
	}

	return nil
}

// Кнопки для настроек
var (
	btnEnableWeather = tele.InlineButton{
		Unique: "enable_weather",
		Text:   "✅ Включить рассылку",
	}
	btnDisableWeather = tele.InlineButton{
		Unique: "disable_weather",
		Text:   "❌ Выключить рассылку",
	}
)

// handleSettings обрабатывает команду /settings
func (s *ApplicationBot) handleSettings(c tele.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	// Получаем текущие настройки пользователя
	user, err := s.storage.GetUser(ctx, chatID)
	if err != nil {
		log.Printf("Failed to get user %d: %v", chatID, err)
		return c.Send("Сначала отправьте команду /start для регистрации.")
	}

	var statusText string
	var keyboard *tele.ReplyMarkup

	if user.WeatherEnabled {
		statusText = "✅ Утренняя рассылка погоды *включена*\n\nВы будете получать прогноз каждый день в 07:00 МСК."
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{btnDisableWeather},
			},
		}
	} else {
		statusText = "❌ Утренняя рассылка погоды *выключена*\n\nВы не будете получать ежедневные прогнозы."
		keyboard = &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{
				{btnEnableWeather},
			},
		}
	}

	return c.Send(statusText, keyboard, tele.ModeMarkdown)
}

// handleEnableWeather обрабатывает нажатие кнопки включения рассылки
func (s *ApplicationBot) handleEnableWeather(c tele.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	if err := s.storage.UpdateWeatherEnabled(ctx, chatID, true); err != nil {
		log.Printf("Failed to enable weather for user %d: %v", chatID, err)
		return c.Respond(
			&tele.CallbackResponse{
				Text: "Произошла ошибка. Попробуйте позже.",
			},
		)
	}

	// Обновляем сообщение
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnDisableWeather},
		},
	}

	if err := c.Edit(
		"✅ Утренняя рассылка погоды *включена*\n\nВы будете получать прогноз каждый день в 07:00 МСК.",
		keyboard,
		tele.ModeMarkdown,
	); err != nil {
		log.Printf("Failed to edit message: %v", err)
	}

	return c.Respond(
		&tele.CallbackResponse{
			Text: "Рассылка включена!",
		},
	)
}

// handleDisableWeather обрабатывает нажатие кнопки выключения рассылки
func (s *ApplicationBot) handleDisableWeather(c tele.Context) error {
	ctx := context.Background()
	chatID := c.Chat().ID

	if err := s.storage.UpdateWeatherEnabled(ctx, chatID, false); err != nil {
		log.Printf("Failed to disable weather for user %d: %v", chatID, err)
		return c.Respond(
			&tele.CallbackResponse{
				Text: "Произошла ошибка. Попробуйте позже.",
			},
		)
	}

	// Обновляем сообщение
	keyboard := &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{btnEnableWeather},
		},
	}

	if err := c.Edit(
		"❌ Утренняя рассылка погоды *выключена*\n\nВы не будете получать ежедневные прогнозы.",
		keyboard,
		tele.ModeMarkdown,
	); err != nil {
		log.Printf("Failed to edit message: %v", err)
	}

	return c.Respond(
		&tele.CallbackResponse{
			Text: "Рассылка выключена.",
		},
	)
}
