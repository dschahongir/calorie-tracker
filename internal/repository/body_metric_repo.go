// internal/repository/body_metric_repo.go
package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"calories-checker/internal/domain"
)

type PostgresBodyMetricRepository struct {
	db *pgxpool.Pool
}

func NewPostgresBodyMetricRepository(db *pgxpool.Pool) *PostgresBodyMetricRepository {
	return &PostgresBodyMetricRepository{db: db}
}

func (r *PostgresBodyMetricRepository) Add(ctx context.Context, metric *domain.BodyMetric) error {
	query := `
		INSERT INTO body_metrics (weight_kg, recorded_at)
		VALUES ($1, $2)
		RETURNING id;
	`
	err := r.db.QueryRow(ctx, query, metric.WeightKG, metric.RecordedAt).Scan(&metric.ID)
	if err != nil {
		return fmt.Errorf("repository: failed to insert body metric: %w", err)
	}
	return nil
}

func (r *PostgresBodyMetricRepository) GetLatest(ctx context.Context) (*domain.BodyMetric, error) {
	query := `
		SELECT id, weight_kg, recorded_at
		FROM body_metrics
		ORDER BY recorded_at DESC
		LIMIT 1;
	`
	var metric domain.BodyMetric
	err := r.db.QueryRow(ctx, query).Scan(&metric.ID, &metric.WeightKG, &metric.RecordedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository: get latest weight: %w", domain.ErrFoodNotFound)
		}
		return nil, fmt.Errorf("repository: failed to fetch latest weight: %w", err)
	}
	return &metric, nil
}
