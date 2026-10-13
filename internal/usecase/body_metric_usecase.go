// internal/usecase/body_metric_usecase.go
package usecase

import (
	"context"
	"fmt"
	"time"

	"calories-checker/internal/domain"
)

type bodyMetricUsecase struct {
	metricRepo domain.BodyMetricRepository
}

func NewBodyMetricUsecase(metricRepo domain.BodyMetricRepository) domain.BodyMetricUsecase {
	return &bodyMetricUsecase{metricRepo: metricRepo}
}

func (u *bodyMetricUsecase) LogWeight(ctx context.Context, weightKG float64) error {
	metric := &domain.BodyMetric{
		WeightKG:   weightKG,
		RecordedAt: time.Now(),
	}

	if err := metric.Validate(); err != nil {
		return fmt.Errorf("usecase: validate weight: %w", domain.ErrInvalidInput)
	}

	if err := u.metricRepo.Add(ctx, metric); err != nil {
		return fmt.Errorf("usecase: log weight: %w", err)
	}
	return nil
}

func (u *bodyMetricUsecase) GetCurrentWeight(ctx context.Context) (*domain.BodyMetric, error) {
	metric, err := u.metricRepo.GetLatest(ctx)
	if err != nil {
		return nil, fmt.Errorf("usecase: get current weight: %w", err)
	}
	return metric, nil
}
