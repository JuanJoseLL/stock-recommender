package main

import (
	"log"

	"github.com/JuanJoseLL/stock-recommender/internal/client"
	"github.com/JuanJoseLL/stock-recommender/internal/handler"
	"github.com/JuanJoseLL/stock-recommender/internal/repository"
	"github.com/JuanJoseLL/stock-recommender/internal/service"
	"github.com/JuanJoseLL/stock-recommender/pkg/config"
	"github.com/JuanJoseLL/stock-recommender/pkg/database"
	"github.com/gin-gonic/gin"
)

func main() {
	
	cfg := config.Load()

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	stockRepo := repository.NewStockRepository(db)
	
	apiClient := client.NewHTTPClient(cfg.API.URL, cfg.API.Key)

	stockService := service.NewStockService(stockRepo, apiClient) 

	
	stockHandler := handler.NewStockHandler(stockService)

	
	router := gin.Default()

	
	stockHandler.RegisterRoutes(router)

	
	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

