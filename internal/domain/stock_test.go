package domain_test

import (
	"testing"
	"time"
	
	"github.com/JuanJoseLL/stock-recommender/internal/domain"
)

func TestStock_Creation(t *testing.T) {
	stock := domain.Stock{
		ID:         1,
		Ticker:     "AAPL",
		TargetFrom: "150.00",
		TargetTo:   "200.00",
		Company:    "Apple Inc.",
		Action:     "BUY",
		Brokerage:  "Goldman Sachs",
		RatingFrom: "HOLD",
		RatingTo:   "BUY",
		Time:       time.Now(),
	}

	if stock.Ticker != "AAPL" {
		t.Errorf("Expected ticker to be AAPL, got %s", stock.Ticker)
	}
	if stock.Company != "Apple Inc." {
		t.Errorf("Expected company to be Apple Inc., got %s", stock.Company)
	}
}

func TestAPIResponse_ItemsLength(t *testing.T) {
	stocks := []domain.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple Inc."},
		{ID: 2, Ticker: "GOOGL", Company: "Google"},
	}

	response := domain.APIResponse{
		Items:    stocks,
		NextPage: "page2",
	}

	if len(response.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(response.Items))
	}
	if response.NextPage != "page2" {
		t.Errorf("Expected NextPage to be page2, got %s", response.NextPage)
	}
}

func TestRecommendationResponse_Confidence(t *testing.T) {
	stock := domain.Stock{
		ID:      1,
		Ticker:  "TSLA",
		Company: "Tesla Inc.",
	}

	recommendation := domain.RecommendationResponse{
		RecommendedStock: stock,
		Confidence:       0.85,
		Reason:           "Strong growth potential",
		AnalysisDate:     "2024-01-15",
	}

	if recommendation.Confidence != 0.85 {
		t.Errorf("Expected confidence to be 0.85, got %f", recommendation.Confidence)
	}
	if recommendation.RecommendedStock.Ticker != "TSLA" {
		t.Errorf("Expected ticker to be TSLA, got %s", recommendation.RecommendedStock.Ticker)
	}
}

func TestStock_EmptyValues(t *testing.T) {
	stock := domain.Stock{}

	if stock.ID != 0 {
		t.Errorf("Expected ID to be 0, got %d", stock.ID)
	}
	if stock.Ticker != "" {
		t.Errorf("Expected empty ticker, got %s", stock.Ticker)
	}
}

func TestAPIResponse_EmptyItems(t *testing.T) {
	response := domain.APIResponse{
		Items:    []domain.Stock{},
		NextPage: "",
	}

	if len(response.Items) != 0 {
		t.Errorf("Expected 0 items, got %d", len(response.Items))
	}
	if response.NextPage != "" {
		t.Errorf("Expected empty NextPage, got %s", response.NextPage)
	}
}