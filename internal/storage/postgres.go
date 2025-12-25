package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// UserRepository определяет интерфейс для работы с пользователями
type UserRepository interface {
	CreateUser(ctx context.Context, chatID int64) error
	GetUser(ctx context.Context, chatID int64) (*User, error)
	UpdateWeatherEnabled(ctx context.Context, chatID int64, enabled bool) error
	GetAllEnabledUsers(ctx context.Context) ([]*User, error)
}

// PostgresStorage реализует UserRepository для PostgreSQL
type PostgresStorage struct {
	db *gorm.DB
}

// NewPostgresStorage создает новое подключение к PostgreSQL и выполняет миграции
func NewPostgresStorage(ctx context.Context, databaseURL string) (*PostgresStorage, error) {
	db, err := gorm.Open(
		postgres.Open(databaseURL), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
			NowFunc: func() time.Time {
				return time.Now().UTC()
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Проверка соединения
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctxWithTimeout); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// AutoMigrate для создания таблиц
	if err := db.AutoMigrate(&User{}); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

// Close закрывает соединение с базой данных
func (s *PostgresStorage) Close() error {
	sqlDB, err := s.db.DB()

	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	return sqlDB.Close()
}

// CreateUser создает нового пользователя или обновляет существующего
func (s *PostgresStorage) CreateUser(ctx context.Context, chatID int64) error {
	user := &User{
		ChatID:         chatID,
		WeatherEnabled: true,
	}

	result := s.db.WithContext(ctx).Where("chat_id = ?", chatID).FirstOrCreate(user)
	if result.Error != nil {
		return fmt.Errorf("failed to create user: %w", result.Error)
	}

	return nil
}

// GetUser получает пользователя по chatID
func (s *PostgresStorage) GetUser(ctx context.Context, chatID int64) (*User, error) {
	var user User

	result := s.db.WithContext(ctx).Where("chat_id = ?", chatID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with chat_id %d not found", chatID)
		}
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

// UpdateWeatherEnabled обновляет статус рассылки погоды для пользователя
func (s *PostgresStorage) UpdateWeatherEnabled(ctx context.Context, chatID int64, enabled bool) error {
	result := s.db.WithContext(ctx).
		Model(&User{}).
		Where("chat_id = ?", chatID).
		Update("weather_enabled", enabled)

	if result.Error != nil {
		return fmt.Errorf("failed to update weather enabled: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with chat_id %d not found", chatID)
	}

	return nil
}

// GetAllEnabledUsers получает всех пользователей с включенной рассылкой погоды
func (s *PostgresStorage) GetAllEnabledUsers(ctx context.Context) ([]*User, error) {
	var users []*User

	result := s.db.WithContext(ctx).Where("weather_enabled = ?", true).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get enabled users: %w", result.Error)
	}

	return users, nil
}

