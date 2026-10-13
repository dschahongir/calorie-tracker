// internal/repository/water_repo.go
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"calories-checker/internal/domain"
)

type PostgresWaterRepository struct {
	db *pgxpool.Pool
}

func NewPostgresWaterRepository(db *pgxpool.Pool) *PostgresWaterRepository {
	return &PostgresWaterRepository{db: db}
}

func (r *PostgresWaterRepository) Add(ctx context.Context, entry *domain.WaterEntry) error {
	query := `
		INSERT INTO water_entries (amount_ml, consumed_at)
		VALUES ($1, $2)
		RETURNING id;
	`
	err := r.db.QueryRow(ctx, query, entry.AmountML, entry.ConsumedAt).Scan(&entry.ID)
	if err != nil {
		return fmt.Errorf("repository: failed to log water entry: %w", err)
	}
	return nil
}

func (r *PostgresWaterRepository) GetTotalByDate(ctx context.Context, date time.Time) (int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	query := `
		SELECT COALESCE(SUM(amount_ml), 0)
		FROM water_entries
		WHERE consumed_at BETWEEN $1 AND $2;
	`
	var total int
	err := r.db.QueryRow(ctx, query, startOfDay, endOfDay).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("repository: failed to calculate total water: %w", err)
	}
	return total, nil
}
