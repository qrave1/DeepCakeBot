package usecase

import (
	"context"
	"log"
	"time"

	"github.com/qrave1/DeepCakeBot/internal/storage"
)

// Scheduler управляет планированием задач
type Scheduler struct {
	storage        storage.UserRepository
	applicationBot *ApplicationBot
	timezone       *time.Location
	scheduleHour   int
	stopChan       chan struct{}
}

// NewScheduler создает новый планировщик
func NewScheduler(storage storage.UserRepository, applicationBot *ApplicationBot, timezoneName string, scheduleHour int) (
	*Scheduler,
	error,
) {
	location, err := time.LoadLocation(timezoneName)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		storage:        storage,
		applicationBot: applicationBot,
		timezone:       location,
		scheduleHour:   scheduleHour,
		stopChan:       make(chan struct{}),
	}, nil
}

// Start запускает планировщик
func (s *Scheduler) Start(ctx context.Context) {
	log.Printf("Scheduler started. Weather will be sent daily at %02d:00 %s", s.scheduleHour, s.timezone.String())

	// Запускаем первую проверку
	go s.run(ctx)
}

// Stop останавливает планировщик
func (s *Scheduler) Stop() {
	close(s.stopChan)
	log.Println("Scheduler stopped")
}

// run основной цикл планировщика
func (s *Scheduler) run(ctx context.Context) {
	nextRun := s.getNextRunTime()
	log.Printf("Next weather broadcast scheduled at: %s", nextRun.Format("2006-01-02 15:04:05 MST"))

	timer := time.NewTimer(time.Until(nextRun))
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler context cancelled")
			return
		case <-s.stopChan:
			log.Println("Scheduler stop signal received")
			return
		case <-timer.C:
			log.Println("Starting scheduled weather broadcast...")
			s.sendWeatherToAllUsers(ctx)

			// Планируем следующий запуск (через 24 часа)
			nextRun = s.getNextRunTime()
			log.Printf("Next weather broadcast scheduled at: %s", nextRun.Format("2006-01-02 15:04:05 MST"))
			timer.Reset(time.Until(nextRun))
		}
	}
}

// getNextRunTime вычисляет время следующего запуска
func (s *Scheduler) getNextRunTime() time.Time {
	now := time.Now().In(s.timezone)

	next := time.Date(
		now.Year(), now.Month(), now.Day(),
		s.scheduleHour, 0, 0, 0,
		s.timezone,
	)

	if now.After(next) || now.Equal(next) {
		next = next.Add(24 * time.Hour)
	}

	return next
}

// sendWeatherToAllUsers отправляет прогноз погоды всем пользователям с включенной рассылкой
func (s *Scheduler) sendWeatherToAllUsers(ctx context.Context) {
	users, err := s.storage.GetAllEnabledUsers(ctx)
	if err != nil {
		log.Printf("Failed to get enabled users: %v", err)
		return
	}

	log.Printf("Sending weather to %d users...", len(users))

	successCount := 0
	failCount := 0

	for _, user := range users {
		if err := s.applicationBot.SendWeatherToUser(ctx, user.ChatID); err != nil {
			log.Printf("Failed to send weather to user %d: %v", user.ChatID, err)
			failCount++
		} else {
			successCount++
		}

		// Небольшая задержка между отправками, чтобы не превысить лимиты Telegram API
		time.Sleep(50 * time.Millisecond)
	}

	log.Printf("Weather broadcast completed. Success: %d, Failed: %d", successCount, failCount)
}

