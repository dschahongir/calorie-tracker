// internal/delivery/http/food_handler.go
package http

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"calories-checker/internal/domain"
)

type FoodHandler struct {
	foodUsecase domain.FoodUsecase
	aiService   domain.AIService
}

func NewFoodHandler(foodUsecase domain.FoodUsecase, aiService domain.AIService) *FoodHandler {
	return &FoodHandler{
		foodUsecase: foodUsecase,
		aiService:   aiService,
	}
}

// DTO для ручного добавления еды
type createFoodRequest struct {
	Name        string    `json:"name" binding:"required"`
	Calories    float64   `json:"calories" binding:"gte=0"`
	Protein     float64   `json:"protein" binding:"gte=0"`
	Carbs       float64   `json:"carbs" binding:"gte=0"`
	Fat         float64   `json:"fat" binding:"gte=0"`
	WeightGrams float64   `json:"weight_grams" binding:"gte=0"`
	ConsumedAt  time.Time `json:"consumed_at"`
}

// RegisterRoutes регистрирует роуты еды в Gin router
func (h *FoodHandler) RegisterRoutes(rg *gin.RouterGroup) {
	foodGroup := rg.Group("/food")
	{
		foodGroup.POST("", h.AddFoodManual)
		foodGroup.POST("/analyze-image", h.AnalyzeImage)
		foodGroup.GET("/daily", h.GetDailySummary)
	}
}

// AddFoodManual обрабатывает ручное добавление блюда
func (h *FoodHandler) AddFoodManual(c *gin.Context) {
	var req createFoodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON body: " + err.Error()})
		return
	}

	entry := &domain.FoodEntry{
		Name:        req.Name,
		Calories:    req.Calories,
		Protein:     req.Protein,
		Carbs:       req.Carbs,
		Fat:         req.Fat,
		WeightGrams: req.WeightGrams,
		IsFromImage: false,
		ConsumedAt:  req.ConsumedAt,
	}

	if err := h.foodUsecase.AddFood(c.Request.Context(), entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// AnalyzeImage принимает multipart/form-data файл (изображение), шлет в Gemini и сохраняет запись
func (h *FoodHandler) AnalyzeImage(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required in 'image' form key"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read uploaded file"})
		return
	}

	mimeType := fileHeader.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	// Передаем частые блюда для повышения точности Gemini (RAG context)
	userFoodContext := []string{"Ош", "Курутоб", "Манту", "Шашлык", "Самбуса", "Лагман"}

	recognized, err := h.aiService.AnalyzeFoodImage(c.Request.Context(), imageBytes, mimeType, userFoodContext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("gemini error: %v", err)})
		return
	}

	entry := &domain.FoodEntry{
		Name:        recognized.Name,
		Calories:    recognized.Calories,
		Protein:     recognized.Protein,
		Carbs:       recognized.Carbs,
		Fat:         recognized.Fat,
		WeightGrams: recognized.WeightGrams,
		IsFromImage: true,
		ConsumedAt:  time.Now(),
	}

	if err := h.foodUsecase.AddFood(c.Request.Context(), entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save recognized food: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

// GetDailySummary выдает статистику за день (по умолчанию сегодня)
func (h *FoodHandler) GetDailySummary(c *gin.Context) {
	dateStr := c.Query("date")
	targetDate := time.Now()

	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
			return
		}
		targetDate = parsedDate
	}

	summary, err := h.foodUsecase.GetDailySummary(c.Request.Context(), targetDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, summary)
}
