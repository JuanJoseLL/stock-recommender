package main

import (
	"log"

	"github.com/JuanJoseLL/stock-recommender/internal/client"
	"github.com/JuanJoseLL/stock-recommender/internal/handler"
	"github.com/JuanJoseLL/stock-recommender/internal/service"
	"github.com/JuanJoseLL/stock-recommender/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	
	cfg := config.Load()

	
	apiClient := client.NewHTTPClient(cfg.API.URL, cfg.API.Key)

	// TODO: Initialize database repository when implemented
	// For now, we'll use a mock repository
	//var stockRepo interface{} = nil

	stockService := service.NewStockService(nil, apiClient) 

	
	stockHandler := handler.NewStockHandler(stockService)

	
	router := gin.Default()

	
	stockHandler.RegisterRoutes(router)

	
	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

