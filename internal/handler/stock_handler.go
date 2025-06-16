package handler

import (
    "net/http"
    "github.com/JuanJoseLL/stock-recommender/internal/service"
    
    "github.com/gin-gonic/gin"
)

// StockHandler handles HTTP requests
type StockHandler struct {
    service *service.StockService
}

// NewStockHandler creates a new handler
func NewStockHandler(service *service.StockService) *StockHandler {
    return &StockHandler{service: service}
}

// RegisterRoutes registers all routes
func (h *StockHandler) RegisterRoutes(router *gin.Engine) {
    api := router.Group("/api")
    {
        api.GET("/stocks", h.GetStocks)
        api.POST("/stocks/sync", h.SyncStocks)
        api.GET("/recommendations", h.GetRecommendations)
    }
    
    router.GET("/health", h.HealthCheck)
}

// GetStocks handles GET /api/stocks
func (h *StockHandler) GetStocks(c *gin.Context) {
    ctx := c.Request.Context()
    
    stocks, err := h.service.GetAllStocks(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve stocks",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "stocks": stocks,
        "count":  len(stocks),
    })
}

// SyncStocks handles POST /api/stocks/sync
func (h *StockHandler) SyncStocks(c *gin.Context) {
    ctx := c.Request.Context()
    
    err := h.service.SyncStocks(ctx)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to sync stocks",
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Stocks synced successfully",
    })
}

// GetRecommendations handles GET /api/recommendations
func (h *StockHandler) GetRecommendations(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "recommendations": []string{},
        "message": "Recommendations feature coming soon",
    })
}

// HealthCheck handles GET /health
func (h *StockHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "service": "stock-recommender",
    })
}