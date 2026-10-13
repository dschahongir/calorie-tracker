// internal/domain/ai.go
package domain

import (
	"context"
)

// RecognizedFood - результат анализа фото от Gemini
type RecognizedFood struct {
	Name        string  `json:"name"`
	Calories    float64 `json:"calories"`
	Protein     float64 `json:"protein"`
	Carbs       float64 `json:"carbs"`
	Fat         float64 `json:"fat"`
	WeightGrams float64 `json:"weight_grams"`
}

// AIService - интерфейс для взаимодействия с Gemini API
type AIService interface {
	AnalyzeFoodImage(ctx context.Context, imageBytes []byte, mimeType string, userFoodContext []string) (*RecognizedFood, error)
}
