package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/JuanJoseLL/stock-recommender/internal/service"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestNewStockHandler(t *testing.T) {
	var stockService *service.StockService = nil
	var recommendationService *service.RecommendationService = nil
	var enrichmentService *service.EnrichmentService = nil
	
	handler := NewStockHandler(stockService, recommendationService, enrichmentService)
	
	if handler == nil {
		t.Error("Expected handler to be created, got nil")
	}
}

func TestStockHandler_HealthCheck(t *testing.T) {
	var stockService *service.StockService = nil
	var recommendationService *service.RecommendationService = nil
	var enrichmentService *service.EnrichmentService = nil
	handler := NewStockHandler(stockService, recommendationService, enrichmentService)
	
	router := setupRouter()
	router.GET("/health", handler.HealthCheck)
	
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response["status"] != "healthy" {
		t.Errorf("Expected status to be healthy, got %v", response["status"])
	}
	
	if response["service"] != "stock-recommender" {
		t.Errorf("Expected service to be stock-recommender, got %v", response["service"])
	}
}

func TestStockHandler_RegisterRoutes(t *testing.T) {
	var stockService *service.StockService = nil
	var recommendationService *service.RecommendationService = nil
	var enrichmentService *service.EnrichmentService = nil
	handler := NewStockHandler(stockService, recommendationService, enrichmentService)
	
	router := gin.New()
	handler.RegisterRoutes(router)
	
	// Test that routes are registered by checking if they exist
	routes := router.Routes()
	
	expectedRoutes := []string{
		"/api/stocks",
		"/api/stocks/sync",
		"/api/stocks/enrich",
		"/api/recommendations",
		"/health",
	}
	
	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Path] = true
	}
	
	for _, expectedRoute := range expectedRoutes {
		if !routeMap[expectedRoute] {
			t.Errorf("Expected route %s to be registered", expectedRoute)
		}
	}
}