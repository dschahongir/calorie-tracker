// internal/usecase/water_usecase.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"calories-checker/internal/domain"
)

type waterUsecase struct {
	waterRepo domain.WaterRepository
}

func NewWaterUsecase(waterRepo domain.WaterRepository) domain.WaterUsecase {
	return &waterUsecase{waterRepo: waterRepo}
}

func (u *waterUsecase) AddWater(ctx context.Context, amountML int) error {
	entry := &domain.WaterEntry{
		AmountML:   amountML,
		ConsumedAt: time.Now(),
	}

	if err := entry.Validate(); err != nil {
		return fmt.Errorf("usecase: validate water: %w", domain.ErrInvalidInput)
	}

	if err := u.waterRepo.Add(ctx, entry); err != nil {
		return fmt.Errorf("usecase: add water: %w", err)
	}
	return nil
}

func (u *waterUsecase) GetDailyWater(ctx context.Context, date time.Time) (int, error) {
	total, err := u.waterRepo.GetTotalByDate(ctx, date)
	if err != nil {
		return 0, fmt.Errorf("usecase: get daily water: %w", err)
	}
	return total, nil
}
