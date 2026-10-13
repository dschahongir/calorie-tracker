// internal/usecase/food_usecase.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"calories-checker/internal/domain"
)

type foodUsecase struct {
	foodRepo domain.FoodRepository
}

// NewFoodUsecase создает экземпляр usecase. Принимает интерфейс domain.FoodRepository
func NewFoodUsecase(foodRepo domain.FoodRepository) domain.FoodUsecase {
	return &foodUsecase{
		foodRepo: foodRepo,
	}
}

// AddFood валидирует запись и сохраняет ее через репозиторий
func (u *foodUsecase) AddFood(ctx context.Context, entry *domain.FoodEntry) error {
	if err := entry.Validate(); err != nil {
		return fmt.Errorf("usecase: validate food entry: %w", domain.ErrInvalidInput)
	}

	if entry.ConsumedAt.IsZero() {
		entry.ConsumedAt = time.Now()
	}
	entry.CreatedAt = time.Now()

	if err := u.foodRepo.Create(ctx, entry); err != nil {
		return fmt.Errorf("usecase: add food entry: %w", err)
	}

	return nil
}

// GetFood получает конкретную запись по ID
func (u *foodUsecase) GetFood(ctx context.Context, id int64) (*domain.FoodEntry, error) {
	if id <= 0 {
		return nil, fmt.Errorf("usecase: invalid food id: %w", domain.ErrInvalidInput)
	}

	entry, err := u.foodRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("usecase: get food: %w", err)
	}

	return entry, nil
}

// GetDailySummary запрашивает все записи за день и высчитывает суммарную энергетическую ценность
func (u *foodUsecase) GetDailySummary(ctx context.Context, date time.Time) (*domain.DailySummary, error) {
	entries, err := u.foodRepo.ListByDate(ctx, date)
	if err != nil {
		return nil, fmt.Errorf("usecase: get daily food entries: %w", err)
	}

	summary := &domain.DailySummary{
		Date:    date,
		Entries: entries,
	}

	for _, entry := range entries {
		summary.TotalCalories += entry.Calories
		summary.TotalProtein += entry.Protein
		summary.TotalCarbs += entry.Carbs
		summary.TotalFat += entry.Fat
	}

	return summary, nil
}



// internal/usecase/food_usecase.go (дополнение)

type extendedFoodUsecase struct {
	foodRepo domain.FoodRepository
	aiClient domain.AIService
}

// NewExtendedFoodUsecase конструктор с поддержкой AI
func NewExtendedFoodUsecase(foodRepo domain.FoodRepository, aiClient domain.AIService) *extendedFoodUsecase {
	return &extendedFoodUsecase{
		foodRepo: foodRepo,
		aiClient: aiClient,
	}
}

// AnalyzeImageAndSave отправляет фото в Gemini и автоматически сохраняет распознанное блюдо в БД
func (u *extendedFoodUsecase) AnalyzeImageAndSave(ctx context.Context, imageBytes []byte, mimeType string, userContext []string) (*domain.FoodEntry, error) {
	recognized, err := u.aiClient.AnalyzeFoodImage(ctx, imageBytes, mimeType, userContext)
	if err != nil {
		return nil, fmt.Errorf("usecase: failed to analyze image via gemini: %w", err)
	}

	entry := &domain.FoodEntry{
		Name:        recognized.Name,
		Calories:    recognized.Calories,
		Protein:     recognized.Protein,
		Carbs:       recognized.Carbs,
		Fat:         recognized.Fat,
		WeightGrams: recognized.WeightGrams,
		IsFromImage: true,
		ConsumedAt:  time.Now(),
		CreatedAt:   time.Now(),
	}

	if err := entry.Validate(); err != nil {
		return nil, fmt.Errorf("usecase: invalid ai output: %w", err)
	}

	if err := u.foodRepo.Create(ctx, entry); err != nil {
		return nil, fmt.Errorf("usecase: failed to save ai recognized food: %w", err)
	}

	return entry, nil
}
