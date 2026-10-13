// internal/repository/gemini_client.go
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"calories-checker/internal/domain"
)

type GeminiClient struct {
	apiKey string
}

// NewGeminiClient создает клиент для обращения к Gemini API
func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{apiKey: apiKey}
}

// AnalyzeFoodImage отправляет изображение в Gemini 1.5 Flash и возвращает структурированный результат
func (g *GeminiClient) AnalyzeFoodImage(ctx context.Context, imageBytes []byte, mimeType string, userFoodContext []string) (*domain.RecognizedFood, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(g.apiKey))
	if err != nil {
		return nil, fmt.Errorf("gemini: failed to create client: %w", err)
	}
	defer client.Close()

	// Используем быстрое и точное мультимодальное решение
	model := client.GenerativeModel("gemini-1.5-flash")
	model.ResponseMIMEType = "application/json"

	// Формируем динамический системный промпт с контекстом локальной еды
	contextPrompt := ""
	if len(userFoodContext) > 0 {
		contextPrompt = fmt.Sprintf("\nFrequently consumed dishes/ingredients by this user: [%s]. Prioritize matching against these if applicable.", strings.Join(userFoodContext, ", "))
	}

	prompt := fmt.Sprintf(`
Analyze this food image. Estimate its nutritional values as accurately as possible.%s

Return ONLY a JSON object with this exact structure:
{
  "name": "Name of dish (in Russian or Tajik if local, e.g., Ош, Курутоб, Манту)",
  "calories": float (total kcal),
  "protein": float (grams),
  "carbs": float (grams),
  "fat": float (grams),
  "weight_grams": float (estimated portion weight in grams)
}
`, contextPrompt)

	promptPart := genai.Text(prompt)
	imagePart := genai.ImageData(mimeType, imageBytes)

	resp, err := model.GenerateContent(ctx, promptPart, imagePart)
	if err != nil {
		return nil, fmt.Errorf("gemini: failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini: empty response from model")
	}

	// Извлекаем текст ответа
	var responseText string
	for _, part := range resp.Candidates[0].Content.Parts {
		if txt, ok := part.(genai.Text); ok {
			responseText += string(txt)
		}
	}

	var result domain.RecognizedFood
	if err := json.Unmarshal([]byte(responseText), &result); err != nil {
		return nil, fmt.Errorf("gemini: failed to parse JSON response (%s): %w", responseText, err)
	}

	return &result, nil
}
