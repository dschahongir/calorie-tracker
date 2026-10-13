// internal/domain/water.go
package domain

import (
	"context"
	"time"
)

type WaterEntry struct {
	ID        int64
	AmountML  int // Объем в миллилитрах
	ConsumedAt time.Time
}

func (w *WaterEntry) Validate() error {
	if w.AmountML <= 0 {
		return ErrInvalidInput
	}
	return nil
}

type WaterRepository interface {
	Add(ctx context.Context, entry *WaterEntry) error
	GetTotalByDate(ctx context.Context, date time.Time) (int, error)
}

type WaterUsecase interface {
	AddWater(ctx context.Context, amountML int) error
	GetDailyWater(ctx context.Context, date time.Time) (int, error)
}
