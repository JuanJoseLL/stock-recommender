package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/JuanJoseLL/stock-recommender/internal/client"
	"github.com/JuanJoseLL/stock-recommender/internal/domain"
)

// Mock repository for recommendation service testing
type mockRecommendationRepository struct {
	stocks []domain.Stock
	err    error
}

func (m *mockRecommendationRepository) GetAll(ctx context.Context) ([]domain.Stock, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stocks, nil
}

func (m *mockRecommendationRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Stock, error) {
	return nil, nil
}

func (m *mockRecommendationRepository) Create(ctx context.Context, stock *domain.Stock) error {
	return m.err
}

func (m *mockRecommendationRepository) Update(ctx context.Context, stock *domain.Stock) error {
	return m.err
}

func (m *mockRecommendationRepository) BulkCreate(ctx context.Context, stocks []domain.Stock) error {
	return m.err
}

func (m *mockRecommendationRepository) GetTopRecommendations(ctx context.Context, limit int) ([]domain.Stock, error) {
	return m.stocks, m.err
}

func (m *mockRecommendationRepository) UpsertStock(ctx context.Context, stock *domain.Stock) error {
	return m.err
}

func (m *mockRecommendationRepository) CheckDuplicates(ctx context.Context) (int, error) {
	return 0, m.err
}

func TestNewRecommendationService(t *testing.T) {
	repo := &mockRecommendationRepository{}
	alphaClient := client.NewAlphaVantageClient("test-key")

	service := NewRecommendationService(repo, alphaClient)

	if service == nil {
		t.Error("Expected service to be created, got nil")
	}
	if service.stockRepo == nil {
		t.Error("Expected repository to be set")
	}
	if service.alphaVantageClient == nil {
		t.Error("Expected client to be set")
	}
}

func TestRecommendationService_GetStockRecommendations_EmptyStocks(t *testing.T) {
	repo := &mockRecommendationRepository{stocks: []domain.Stock{}}
	alphaClient := client.NewAlphaVantageClient("test-key")
	service := NewRecommendationService(repo, alphaClient)

	response, err := service.GetStockRecommendations(context.Background(), 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == nil {
		t.Error("Expected response, got nil")
		return
	}
	if len(response.Recommendations) != 0 {
		t.Errorf("Expected 0 recommendations, got %d", len(response.Recommendations))
	}
	if response.Summary.TotalAnalyzed != 0 {
		t.Errorf("Expected total analyzed to be 0, got %d", response.Summary.TotalAnalyzed)
	}
}

func TestRecommendationService_GetStockRecommendations_WithStocks(t *testing.T) {
	testStocks := []domain.Stock{
		{
			Ticker:     "AAPL",
			Company:    "Apple Inc.",
			Action:     "BUY",
			Brokerage:  "Goldman Sachs",
			RatingFrom: "HOLD",
			RatingTo:   "BUY",
		},
		{
			Ticker:    "GOOGL",
			Company:   "Google LLC",
			Action:    "HOLD",
			Brokerage: "Morgan Stanley",
		},
	}

	repo := &mockRecommendationRepository{stocks: testStocks}
	alphaClient := client.NewAlphaVantageClient("test-key")
	service := NewRecommendationService(repo, alphaClient)

	response, err := service.GetStockRecommendations(context.Background(), 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if response == nil {
		t.Error("Expected response, got nil")
		return
	}
	if len(response.Recommendations) == 0 {
		t.Error("Expected recommendations, got none")
	}
	if response.Summary.TotalAnalyzed != 2 {
		t.Errorf("Expected total analyzed to be 2, got %d", response.Summary.TotalAnalyzed)
	}
}

func TestRecommendationService_GetStockRecommendations_DatabaseError(t *testing.T) {
	repo := &mockRecommendationRepository{err: errors.New("database error")}
	alphaClient := client.NewAlphaVantageClient("test-key")
	service := NewRecommendationService(repo, alphaClient)

	response, err := service.GetStockRecommendations(context.Background(), 10)

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if response != nil {
		t.Error("Expected nil response on error")
	}
}

func TestRecommendationService_calculateGainerScore(t *testing.T) {
	service := &RecommendationService{}

	tests := []struct {
		name             string
		changePercentage float64
		volume           string
		expectedMinScore float64
		expectedMaxScore float64
	}{
		{
			name:             "High percentage with high volume",
			changePercentage: 5.0,
			volume:           "2000000",
			expectedMinScore: 70, // 5*10 + 20 volume bonus
			expectedMaxScore: 100,
		},
		{
			name:             "Low percentage with low volume",
			changePercentage: 1.0,
			volume:           "50000",
			expectedMinScore: 10, // 1*10
			expectedMaxScore: 20,
		},
		{
			name:             "Very high percentage capped at 100",
			changePercentage: 15.0,
			volume:           "3000000",
			expectedMinScore: 100, // Capped at 100
			expectedMaxScore: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateGainerScore(tt.changePercentage, tt.volume)
			if score < tt.expectedMinScore || score > tt.expectedMaxScore {
				t.Errorf("Expected score between %f and %f, got %f", tt.expectedMinScore, tt.expectedMaxScore, score)
			}
		})
	}
}

func TestRecommendationService_calculateAnalystScore(t *testing.T) {
	service := &RecommendationService{}

	tests := []struct {
		name          string
		stock         domain.Stock
		expectedRange [2]float64 // min, max expected score
	}{
		{
			name: "Strong BUY action",
			stock: domain.Stock{
				Action:    "STRONG BUY",
				Brokerage: "Goldman Sachs",
			},
			expectedRange: [2]float64{45, 55}, // 40 + 10 bonus
		},
		{
			name: "BUY action with rating upgrade",
			stock: domain.Stock{
				Action:     "BUY",
				RatingFrom: "HOLD",
				RatingTo:   "BUY",
				Brokerage:  "Local Firm",
			},
			expectedRange: [2]float64{50, 60}, // 40 + 15 upgrade
		},
		{
			name: "SELL action",
			stock: domain.Stock{
				Action:    "SELL",
				Brokerage: "Some Firm",
			},
			expectedRange: [2]float64{-35, -25}, // -30
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := service.calculateAnalystScore(tt.stock)
			if score < tt.expectedRange[0] || score > tt.expectedRange[1] {
				t.Errorf("Expected score between %f and %f, got %f", tt.expectedRange[0], tt.expectedRange[1], score)
			}
		})
	}
}

func TestRecommendationService_parsePrice(t *testing.T) {
	service := &RecommendationService{}

	tests := []struct {
		priceStr string
		expected float64
	}{
		{"$150.00", 150.00},
		{"1,250.50", 1250.50},
		{"", 0},
		{"invalid", 0},
		{"100", 100},
		{"$1,000.25", 1000.25},
	}

	for _, tt := range tests {
		t.Run(tt.priceStr, func(t *testing.T) {
			result := service.parsePrice(tt.priceStr)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

func TestRecommendationService_getRecommendationType(t *testing.T) {
	service := &RecommendationService{}

	tests := []struct {
		score    float64
		expected string
	}{
		{80, "BUY"},
		{70, "BUY"},
		{60, "HOLD"},
		{50, "HOLD"},
		{40, "WATCH"},
		{0, "WATCH"},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := service.getRecommendationType(tt.score)
			if result != tt.expected {
				t.Errorf("For score %f, expected %s, got %s", tt.score, tt.expected, result)
			}
		})
	}
}

func TestRecommendationService_generateBasicReason(t *testing.T) {
	service := &RecommendationService{}

	stock := domain.Stock{
		Action:     "BUY",
		Brokerage:  "Goldman Sachs",
		RatingFrom: "HOLD",
		RatingTo:   "BUY",
		TargetFrom: "150.00",
		TargetTo:   "200.00",
	}

	reason := service.generateBasicReason(stock, 75.0)

	if !strings.Contains(reason, "buy recommendation") {
		t.Error("Expected reason to contain 'buy recommendation'")
	}
	if !strings.Contains(reason, "Goldman Sachs") {
		t.Error("Expected reason to contain 'Goldman Sachs'")
	}
	if !strings.Contains(reason, "75.0/100") {
		t.Error("Expected reason to contain score")
	}
}

func TestRecommendationService_analyzeTopGainers(t *testing.T) {
	service := &RecommendationService{}

	gainers := []client.AlphaVantageMarketMover{
		{
			Ticker:           "TSLA",
			Price:            "250.00",
			ChangeAmount:     "12.50",
			ChangePercentage: "5.26%",
			Volume:           "2500000",
		},
		{
			Ticker:           "NVDA",
			Price:            "400.00",
			ChangeAmount:     "20.00",
			ChangePercentage: "8.33%",
			Volume:           "1500000",
		},
	}

	recommendations := service.analyzeTopGainers(gainers)

	if len(recommendations) != 2 {
		t.Errorf("Expected 2 recommendations, got %d", len(recommendations))
	}

	if recommendations[0].Symbol != "TSLA" {
		t.Errorf("Expected first recommendation to be TSLA, got %s", recommendations[0].Symbol)
	}

	if recommendations[0].Score <= 0 {
		t.Error("Expected positive score for top gainer")
	}

	if recommendations[0].RecommendationType == "" {
		t.Error("Expected recommendation type to be set")
	}
}

func TestRecommendationService_analyzeStoredStocks(t *testing.T) {
	service := &RecommendationService{}

	stocks := []domain.Stock{
		{
			Ticker:    "AAPL",
			Company:   "Apple Inc.",
			Action:    "BUY",
			Brokerage: "Goldman Sachs",
		},
		{
			Ticker:    "GOOGL",
			Company:   "Google LLC",
			Action:    "HOLD",
			Brokerage: "Morgan Stanley",
		},
	}

	recommendations := service.analyzeStoredStocks(stocks)

	if len(recommendations) != 2 {
		t.Errorf("Expected 2 recommendations, got %d", len(recommendations))
	}

	if recommendations[0].Symbol != "AAPL" {
		t.Errorf("Expected first recommendation to be AAPL, got %s", recommendations[0].Symbol)
	}

	if recommendations[0].Name != "Apple Inc." {
		t.Errorf("Expected first recommendation name to be Apple Inc., got %s", recommendations[0].Name)
	}
}
