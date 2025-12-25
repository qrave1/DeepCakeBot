package storage

import "gorm.io/gorm"

// User представляет пользователя бота
type User struct {
	gorm.Model
	// ChatID - уникальный идентификатор чата в Telegram
	ChatID int64 `gorm:"uniqueIndex;not null"`
	// WeatherEnabled - флаг включения утренней рассылки погоды
	WeatherEnabled bool `gorm:"default:true;not null"`
}

