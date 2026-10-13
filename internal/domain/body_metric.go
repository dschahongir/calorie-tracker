// internal/domain/body_metric.go
package domain

import (
	"context"
	"time"
)

type BodyMetric struct {
	ID        int64
	WeightKG  float64
	RecordedAt time.Time
}

func (b *BodyMetric) Validate() error {
	if b.WeightKG <= 0 {
		return ErrInvalidInput
	}
	return nil
}

type BodyMetricRepository interface {
	Add(ctx context.Context, metric *BodyMetric) error
	GetLatest(ctx context.Context) (*BodyMetric, error)
}

type BodyMetricUsecase interface {
	LogWeight(ctx context.Context, weightKG float64) error
	GetCurrentWeight(ctx context.Context) (*BodyMetric, error)
}
