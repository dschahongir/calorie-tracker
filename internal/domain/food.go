// internal/domain/food.go
package domain

import (
	"context"
	"errors"
	"time"
)

// Доменные ошибки
var (
	ErrFoodNotFound = errors.New("food not found")
	ErrInvalidInput = errors.New("invalid input data")
)

// FoodEntry представляет запись о съеденной пище (добавлено вручную или через фото)
type FoodEntry struct {
	ID            int64
	Name          string
	Calories      float64
	Protein       float64
	Carbs         float64
	Fat           float64
	WeightGrams   float64 // Вес порции в граммах
	IsFromImage   bool    // Флаг, была ли еда добавлена через фото/Gemini
	ImageURL      *string // Опционально ссылка на загруженное фото
	ConsumedAt    time.Time
	CreatedAt     time.Time
}

// Validation бизнес-правила внутри entity
func (f *FoodEntry) Validate() error {
	if f.Name == "" {
		return errors.New("food name is required")
	}
	if f.Calories < 0 || f.Protein < 0 || f.Carbs < 0 || f.Fat < 0 || f.WeightGrams < 0 {
		return errors.New("nutritional values and weight cannot be negative")
	}
	return nil
}

// Интерфейс репозитория (для CRUD в БД)
type FoodRepository interface {
	Create(ctx context.Context, entry *FoodEntry) error
	GetByID(ctx context.Context, id int64) (*FoodEntry, error)
	ListByDate(ctx context.Context, date time.Time) ([]*FoodEntry, error)
}

// Интерфейс бизнес-логики (usecase)
type FoodUsecase interface {
	AddFood(ctx context.Context, entry *FoodEntry) error
	GetFood(ctx context.Context, id int64) (*FoodEntry, error)
	GetDailySummary(ctx context.Context, date time.Time) (*DailySummary, error)
}

// DailySummary - DTO для красивых графиков и статистики за день
type DailySummary struct {
	Date          time.Time
	TotalCalories float64
	TotalProtein  float64
	TotalCarbs    float64
	TotalFat      float64
	Entries       []*FoodEntry
}
