package main

import (
	"log"

	"time"

	"github.com/JuanJoseLL/stock-recommender/internal/client"
	"github.com/JuanJoseLL/stock-recommender/internal/handler"
	"github.com/JuanJoseLL/stock-recommender/internal/repository"
	"github.com/JuanJoseLL/stock-recommender/internal/service"
	"github.com/JuanJoseLL/stock-recommender/pkg/config"
	"github.com/JuanJoseLL/stock-recommender/pkg/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	// Initialize database connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repository
	stockRepo := repository.NewStockRepository(db)

	// Initialize clients
	apiClient := client.NewHTTPClient(cfg.API.URL, cfg.API.Key)
	alphaVantageClient := client.NewAlphaVantageClient(cfg.AlphaVantage.APIKey)

	// Initialize services
	stockService := service.NewStockService(stockRepo, apiClient)
	recommendationService := service.NewRecommendationService(stockRepo, alphaVantageClient)
	enrichmentService := service.NewEnrichmentService(stockRepo, alphaVantageClient)

	// Initialize handler
	stockHandler := handler.NewStockHandler(stockService, recommendationService, enrichmentService)

	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // En producción, especifica el dominio de CloudFront
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	stockHandler.RegisterRoutes(router)

	log.Printf("Starting server on port %s", cfg.Server.Port)
	log.Printf("Database URL: %s", cfg.GetDatabaseURL())
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
