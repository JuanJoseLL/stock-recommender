package handler

import (
    "net/http"
    "strconv"
    "github.com/JuanJoseLL/stock-recommender/internal/service"
    
    "github.com/gin-gonic/gin"
)

// StockHandler handles HTTP requests
type StockHandler struct {
    service *service.StockService
    recommendationService *service.RecommendationService
    enrichmentService *service.EnrichmentService
}

// NewStockHandler creates a new handler
func NewStockHandler(service *service.StockService, recommendationService *service.RecommendationService, enrichmentService *service.EnrichmentService) *StockHandler {
    return &StockHandler{
        service: service,
        recommendationService: recommendationService,
        enrichmentService: enrichmentService,
    }
}

// RegisterRoutes registers all routes
func (h *StockHandler) RegisterRoutes(router *gin.Engine) {
    api := router.Group("/api")
    {
        api.GET("/stocks", h.GetStocks)
        api.POST("/stocks/sync", h.SyncStocks)
        api.POST("/stocks/enrich", h.EnrichStocks)
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
    ctx := c.Request.Context()
    
    // Get limit parameter (default to 10)
    limitStr := c.DefaultQuery("limit", "10")
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 10
    }
    
    if limit > 50 {
        limit = 50 // Cap at 50 recommendations
    }
    
    recommendations, err := h.recommendationService.GetStockRecommendations(ctx, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to get recommendations",
        })
        return
    }
    
    c.JSON(http.StatusOK, recommendations)
}

// HealthCheck handles GET /health
func (h *StockHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "service": "stock-recommender",
    })
}

// EnrichStocks handles POST /api/stocks/enrich
func (h *StockHandler) EnrichStocks(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Get limit parameter (default to 5 for rate limiting)
    limitStr := c.DefaultQuery("limit", "5")
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 {
        limit = 5
    }
    
    if limit > 20 {
        limit = 20 // Cap at 20 to respect Alpha Vantage rate limits
    }
    
    stats, err := h.enrichmentService.EnrichStockData(ctx, limit)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to enrich stock data",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Stock enrichment completed",
        "stats": stats,
    })
}