package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qrave1/DeepCakeBot/internal/config"
	"github.com/qrave1/DeepCakeBot/internal/storage"
	"github.com/qrave1/DeepCakeBot/internal/usecase"

	tele "gopkg.in/telebot.v3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println("Starting DeepCake Bot...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := storage.NewPostgresStorage(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()
	log.Println("Connected to database")

	pref := tele.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(pref)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Println("Bot created successfully")

	weatherService := usecase.NewWeatherService(cfg.OpenWeatherAPIKey, cfg.City, cfg.CountryCode)

	applicationBot := usecase.NewApplicationBot(bot, db, weatherService)

	applicationBot.RegisterHandlers()
	log.Println("Bot handlers registered")

	scheduler, err := usecase.NewScheduler(db, applicationBot, cfg.Timezone, cfg.WeatherScheduleHour)
	if err != nil {
		log.Fatalf("Failed to create scheduler: %v", err)
	}
	scheduler.Start(ctx)

	go bot.Start()
	log.Println("Bot started and listening for messages...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-sigChan
	log.Println("Shutdown signal received, gracefully shutting down...")

	scheduler.Stop()

	bot.Stop()

	cancel()

	time.Sleep(2 * time.Second)

	log.Println("Bot stopped successfully")
}
