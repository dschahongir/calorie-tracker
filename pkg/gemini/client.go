// pkg/gemini/client.go
package gemini

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"calories-checker/internal/domain"
)

// Client реализует интерфейс domain.AIService
type Client struct {
	apiKey string
}

// NewClient создает новый экземпляр клиента Gemini
func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

// AnalyzeFoodImage отправляет фото в Gemini и возвращает разобранный JSON
func (c *Client) AnalyzeFoodImage(ctx context.Context, imageBytes []byte, mimeType string, userFoodContext []string) (*domain.RecognizedFood, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(c.apiKey))
	if err != nil {
		return nil, fmt.Errorf("gemini: failed to create client: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.5-flash-lite")
	model.ResponseMIMEType = "application/json"

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
