// internal/repository/food_repo.go
package repository

import (
	"context"
	"fmt"
	"errors"
	"time"
	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"

	"calories-checker/internal/domain"
)

// PostgresFoodRepository реализует domain.FoodRepository
type PostgresFoodRepository struct {
	db *pgxpool.Pool // Приватное поле, доступ только внутри пакета repository
}

// NewPostgresFoodRepository — конструктор репозитория
func NewPostgresFoodRepository(db *pgxpool.Pool) *PostgresFoodRepository {
	return &PostgresFoodRepository{db: db}
}

// Create вставляет новую запись о еде в БД и возвращает полученный ID
func (r *PostgresFoodRepository) Create(ctx context.Context, entry *domain.FoodEntry) error {
	query := `
		INSERT INTO food_entries (
			name, calories, protein, carbs, fat, weight_grams, is_from_image, image_url, consumed_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		entry.Name,
		entry.Calories,
		entry.Protein,
		entry.Carbs,
		entry.Fat,
		entry.WeightGrams,
		entry.IsFromImage,
		entry.ImageURL,
		entry.ConsumedAt,
		entry.CreatedAt,
	).Scan(&entry.ID)

	if err != nil {
		return fmt.Errorf("repository: failed to insert food entry: %w", err)
	}

	return nil
}


// GetByID возвращает запись о приеме пищи по её ID
func (r *PostgresFoodRepository) GetByID(ctx context.Context, id int64) (*domain.FoodEntry, error) {
	query := `
		SELECT id, name, calories, protein, carbs, fat, weight_grams, is_from_image, image_url, consumed_at, created_at
		FROM food_entries
		WHERE id = $1;
	`

	var entry domain.FoodEntry
	err := r.db.QueryRow(ctx, query, id).Scan(
		&entry.ID,
		&entry.Name,
		&entry.Calories,
		&entry.Protein,
		&entry.Carbs,
		&entry.Fat,
		&entry.WeightGrams,
		&entry.IsFromImage,
		&entry.ImageURL,
		&entry.ConsumedAt,
		&entry.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository: get food by id: %w", domain.ErrFoodNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get food entry: %w", err)
	}

	return &entry, nil
}

// ListByDate возвращает список всех записей за конкретную дату (от 00:00:00 до 23:59:59)
func (r *PostgresFoodRepository) ListByDate(ctx context.Context, date time.Time) ([]*domain.FoodEntry, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24*time.Hour - time.Nanosecond)

	query := `
		SELECT id, name, calories, protein, carbs, fat, weight_grams, is_from_image, image_url, consumed_at, created_at
		FROM food_entries
		WHERE consumed_at BETWEEN $1 AND $2
		ORDER BY consumed_at ASC;
	`

	rows, err := r.db.Query(ctx, query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to query food entries by date: %w", err)
	}
	defer rows.Close()

	var entries []*domain.FoodEntry
	for rows.Next() {
		var entry domain.FoodEntry
		if err := rows.Scan(
			&entry.ID,
			&entry.Name,
			&entry.Calories,
			&entry.Protein,
			&entry.Carbs,
			&entry.Fat,
			&entry.WeightGrams,
			&entry.IsFromImage,
			&entry.ImageURL,
			&entry.ConsumedAt,
			&entry.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repository: failed to scan food entry row: %w", err)
		}
		entries = append(entries, &entry)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows iteration error: %w", err)
	}

	return entries, nil
}
