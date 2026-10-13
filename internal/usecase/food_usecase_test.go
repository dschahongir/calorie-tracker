// internal/usecase/food_usecase_test.go
package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"calories-checker/internal/domain"
	"calories-checker/internal/usecase"
)

// MockFoodRepository имитирует domain.FoodRepository
type MockFoodRepository struct {
	mock.Mock
}

func (m *MockFoodRepository) Create(ctx context.Context, entry *domain.FoodEntry) error {
	args := m.Called(ctx, entry)
	if args.Get(0) != nil {
		entry.ID = 1
	}
	return args.Error(0)
}

func (m *MockFoodRepository) GetByID(ctx context.Context, id int64) (*domain.FoodEntry, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.FoodEntry), args.Error(1)
}

func (m *MockFoodRepository) ListByDate(ctx context.Context, date time.Time) ([]*domain.FoodEntry, error) {
	args := m.Called(ctx, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.FoodEntry), args.Error(1)
}

func TestFoodUsecase_AddFood_Success(t *testing.T) {
	mockRepo := new(MockFoodRepository)
	uc := usecase.NewFoodUsecase(mockRepo)

	entry := &domain.FoodEntry{
		Name:        "Манту",
		Calories:    450,
		Protein:     18,
		Carbs:       50,
		Fat:         18,
		WeightGrams: 250,
	}

	mockRepo.On("Create", mock.Anything, entry).Return(nil)

	err := uc.AddFood(context.Background(), entry)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), entry.ID)
	mockRepo.AssertExpectations(t)
}

func TestFoodUsecase_AddFood_InvalidInput(t *testing.T) {
	mockRepo := new(MockFoodRepository)
	uc := usecase.NewFoodUsecase(mockRepo)

	// Невалидная сущность (пустое имя)
	entry := &domain.FoodEntry{
		Name:     "",
		Calories: -100,
	}

	err := uc.AddFood(context.Background(), entry)

	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrInvalidInput)
	mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
}

func TestFoodUsecase_GetDailySummary_Success(t *testing.T) {
	mockRepo := new(MockFoodRepository)
	uc := usecase.NewFoodUsecase(mockRepo)

	now := time.Now()
	mockEntries := []*domain.FoodEntry{
		{Name: "Завтрак", Calories: 300, Protein: 10, Carbs: 40, Fat: 10},
		{Name: "Обед", Calories: 600, Protein: 30, Carbs: 70, Fat: 20},
	}

	mockRepo.On("ListByDate", mock.Anything, now).Return(mockEntries, nil)

	summary, err := uc.GetDailySummary(context.Background(), now)

	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 900.0, summary.TotalCalories)
	assert.Equal(t, 40.0, summary.TotalProtein)
	assert.Equal(t, 110.0, summary.TotalCarbs)
	assert.Equal(t, 30.0, summary.TotalFat)
	assert.Len(t, summary.Entries, 2)
	mockRepo.AssertExpectations(t)
}

