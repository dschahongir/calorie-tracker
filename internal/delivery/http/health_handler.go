// internal/delivery/http/health_handler.go
package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"calories-checker/internal/domain"
)

type HealthHandler struct {
	waterUsecase      domain.WaterUsecase
	bodyMetricUsecase domain.BodyMetricUsecase
}

func NewHealthHandler(waterUsecase domain.WaterUsecase, bodyMetricUsecase domain.BodyMetricUsecase) *HealthHandler {
	return &HealthHandler{
		waterUsecase:      waterUsecase,
		bodyMetricUsecase: bodyMetricUsecase,
	}
}

type addWaterRequest struct {
	AmountML int `json:"amount_ml" binding:"required,gt=0"`
}

type logWeightRequest struct {
	WeightKG float64 `json:"weight_kg" binding:"required,gt=0"`
}

func (h *HealthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.POST("/water", h.AddWater)
	rg.GET("/water/daily", h.GetDailyWater)

	rg.POST("/weight", h.LogWeight)
	rg.GET("/weight/latest", h.GetLatestWeight)
}

func (h *HealthHandler) AddWater(c *gin.Context) {
	var req addWaterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid water amount"})
		return
	}

	if err := h.waterUsecase.AddWater(c.Request.Context(), req.AmountML); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "water entry logged successfully"})
}

func (h *HealthHandler) GetDailyWater(c *gin.Context) {
	total, err := h.waterUsecase.GetDailyWater(c.Request.Context(), time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_ml": total})
}

func (h *HealthHandler) LogWeight(c *gin.Context) {
	var req logWeightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid weight value"})
		return
	}

	if err := h.bodyMetricUsecase.LogWeight(c.Request.Context(), req.WeightKG); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "weight logged successfully"})
}

func (h *HealthHandler) GetLatestWeight(c *gin.Context) {
	metric, err := h.bodyMetricUsecase.GetCurrentWeight(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metric)
}
