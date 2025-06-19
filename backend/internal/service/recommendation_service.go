package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/JuanJoseLL/stock-recommender/internal/client"
	"github.com/JuanJoseLL/stock-recommender/internal/domain"
	"github.com/JuanJoseLL/stock-recommender/internal/repository"
)

type RecommendationService struct {
	stockRepo       repository.StockRepository
	alphaVantageClient *client.AlphaVantageClient
}

type StockRecommendation struct {
	Symbol          string  `json:"symbol"`
	Name            string  `json:"name"`
	Score           float64 `json:"score"`
	Reason          string  `json:"reason"`
	CurrentPrice    string  `json:"current_price"`
	TargetPrice     string  `json:"target_price,omitempty"`
	PERatio         string  `json:"pe_ratio,omitempty"`
	DividendYield   string  `json:"dividend_yield,omitempty"`
	MarketCap       string  `json:"market_cap,omitempty"`
	Sector          string  `json:"sector,omitempty"`
	RecommendationType string `json:"recommendation_type"`
}

type RecommendationResponse struct {
	Recommendations []StockRecommendation `json:"recommendations"`
	Summary         RecommendationSummary `json:"summary"`
}

type RecommendationSummary struct {
	TotalAnalyzed    int    `json:"total_analyzed"`
	BuyRecommendations int  `json:"buy_recommendations"`
	HoldRecommendations int `json:"hold_recommendations"`
	GeneratedAt      string `json:"generated_at"`
	DataSource       string `json:"data_source"`
}

func NewRecommendationService(stockRepo repository.StockRepository, alphaVantageClient *client.AlphaVantageClient) *RecommendationService {
	return &RecommendationService{
		stockRepo:       stockRepo,
		alphaVantageClient: alphaVantageClient,
	}
}

func (s *RecommendationService) GetStockRecommendations(ctx context.Context, limit int) (*RecommendationResponse, error) {
	stocks, err := s.stockRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks from database: %w", err)
	}

	if len(stocks) == 0 {
		return &RecommendationResponse{
			Recommendations: []StockRecommendation{},
			Summary: RecommendationSummary{
				TotalAnalyzed: 0,
				BuyRecommendations: 0,
				HoldRecommendations: 0,
				GeneratedAt: "now",
				DataSource: "Alpha Vantage + Database",
			},
		}, nil
	}

	var recommendations []StockRecommendation
	
	topGainers, err := s.alphaVantageClient.GetTopGainersLosers(ctx)
	if err == nil && topGainers != nil {
		recommendations = append(recommendations, s.analyzeTopGainers(topGainers.TopGainers)...)
	}

	dbRecommendations := s.analyzeStoredStocks(stocks)
	recommendations = append(recommendations, dbRecommendations...)

	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	if limit > 0 && len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	buyCount := 0
	holdCount := 0
	for _, rec := range recommendations {
		if rec.RecommendationType == "BUY" {
			buyCount++
		} else if rec.RecommendationType == "HOLD" {
			holdCount++
		}
	}

	return &RecommendationResponse{
		Recommendations: recommendations,
		Summary: RecommendationSummary{
			TotalAnalyzed: len(stocks),
			BuyRecommendations: buyCount,
			HoldRecommendations: holdCount,
			GeneratedAt: "now",
			DataSource: "Alpha Vantage + Database",
		},
	}, nil
}

func (s *RecommendationService) analyzeTopGainers(gainers []client.AlphaVantageMarketMover) []StockRecommendation {
	var recommendations []StockRecommendation
	
	for i, gainer := range gainers {
		if i >= 5 { // Limit to top 5 gainers
			break
		}
		
		changePercentage := strings.TrimSuffix(gainer.ChangePercentage, "%")
		changeFloat, err := strconv.ParseFloat(changePercentage, 64)
		if err != nil {
			continue
		}

		score := s.calculateGainerScore(changeFloat, gainer.Volume)
		
		recommendation := StockRecommendation{
			Symbol: gainer.Ticker,
			Name:   gainer.Ticker, // Alpha Vantage doesn't provide company name in this endpoint
			Score:  score,
			Reason: fmt.Sprintf("Top gainer with %.2f%% increase. High momentum stock with significant volume.", changeFloat),
			CurrentPrice: gainer.Price,
			RecommendationType: s.getRecommendationType(score),
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	return recommendations
}

func (s *RecommendationService) analyzeStoredStocks(stocks []domain.Stock) []StockRecommendation {
	var recommendations []StockRecommendation
	
	for _, stock := range stocks {
		if len(recommendations) >= 10 { // Limit database recommendations
			break
		}
		
		score := s.calculateBasicScore(stock)
		
		recommendation := StockRecommendation{
			Symbol: stock.Ticker,
			Name:   stock.Company,
			Score:  score,
			Reason: s.generateBasicReason(stock, score),
			CurrentPrice: "N/A", // Price not available in current stock model
			RecommendationType: s.getRecommendationType(score),
		}
		
		recommendations = append(recommendations, recommendation)
	}
	
	return recommendations
}

func (s *RecommendationService) calculateGainerScore(changePercentage float64, volumeStr string) float64 {
	volume, err := strconv.ParseFloat(volumeStr, 64)
	if err != nil {
		volume = 0
	}
	
	// Base score from percentage change
	score := changePercentage * 10
	
	// Volume bonus (higher volume = more reliable)
	if volume > 1000000 {
		score += 20
	} else if volume > 500000 {
		score += 10
	} else if volume > 100000 {
		score += 5
	}
	
	// Cap the score at 100
	if score > 100 {
		score = 100
	}
	
	return math.Max(0, score)
}

func (s *RecommendationService) calculateBasicScore(stock domain.Stock) float64 {
	score := 50.0 // Base score
	
	// Analyst sentiment scoring (40% weight)
	analystScore := s.calculateAnalystScore(stock)
	score += analystScore * 0.4
	
	// Fundamental analysis scoring (30% weight)
	fundamentalScore := s.calculateFundamentalScore(stock)
	score += fundamentalScore * 0.3
	
	// Price target analysis scoring (30% weight)
	targetScore := s.calculateTargetScore(stock)
	score += targetScore * 0.3
	
	return math.Max(0, math.Min(100, score))
}

func (s *RecommendationService) calculateAnalystScore(stock domain.Stock) float64 {
	score := 0.0
	
	// Action-based scoring
	switch strings.ToUpper(stock.Action) {
	case "BUY", "STRONG BUY":
		score += 40
	case "OUTPERFORM", "OVERWEIGHT", "UPGRADED BY":
		score += 30
	case "INITIATED BY":
		if strings.Contains(strings.ToUpper(stock.RatingTo), "BUY") {
			score += 35
		} else {
			score += 15
		}
	case "HOLD", "NEUTRAL":
		score += 10
	case "UNDERPERFORM", "UNDERWEIGHT", "DOWNGRADED BY":
		score -= 15
	case "SELL", "STRONG SELL":
		score -= 30
	case "TARGET LOWERED BY", "TARGET RAISED BY":
		score += 5 // Neutral, just analyst attention
	}
	
	// Rating upgrade/downgrade consideration
	if stock.RatingFrom != "" && stock.RatingTo != "" {
		from := strings.ToUpper(stock.RatingFrom)
		to := strings.ToUpper(stock.RatingTo)
		
		if from != to {
			// Upgrade scenarios
			if (from == "HOLD" || from == "NEUTRAL") && strings.Contains(to, "BUY") {
				score += 15
			} else if strings.Contains(from, "SELL") && (strings.Contains(to, "BUY") || to == "HOLD") {
				score += 20
			} else if strings.Contains(from, "BUY") && strings.Contains(to, "SELL") {
				score -= 25 // Significant downgrade
			}
		}
	}
	
	// Brokerage reputation bonus
	if strings.Contains(strings.ToLower(stock.Brokerage), "goldman") ||
		strings.Contains(strings.ToLower(stock.Brokerage), "morgan") ||
		strings.Contains(strings.ToLower(stock.Brokerage), "jp") {
		score += 10
	}
	
	return score
}

func (s *RecommendationService) calculateFundamentalScore(stock domain.Stock) float64 {
	score := 0.0
	
	// PE Ratio analysis
	if stock.PERatio != nil {
		pe := *stock.PERatio
		if pe > 0 {
			if pe < 15 {
				score += 20 // Undervalued
			} else if pe < 25 {
				score += 10 // Reasonable valuation
			} else if pe > 40 {
				score -= 10 // Overvalued
			}
		}
	}
	
	// 52-week range analysis
	if stock.CurrentPrice != nil && stock.WeekHigh52 != nil && stock.WeekLow52 != nil {
		current := *stock.CurrentPrice
		high52 := *stock.WeekHigh52
		low52 := *stock.WeekLow52
		
		if high52 > low52 && current > 0 {
			// Position within 52-week range
			position := (current - low52) / (high52 - low52)
			
			if position < 0.3 {
				score += 15 // Near 52-week low, potential value
			} else if position > 0.8 {
				score -= 5 // Near 52-week high, less upside
			} else {
				score += 5 // Middle range
			}
		}
	}
	
	// Market cap consideration
	if stock.MarketCap != nil {
		marketCap := *stock.MarketCap
		if marketCap > 10000000000 { // Large cap (>10B)
			score += 5 // Stability
		} else if marketCap > 2000000000 { // Mid cap (2B-10B)
			score += 10 // Growth potential
		} else if marketCap > 300000000 { // Small cap (300M-2B)
			score += 8 // Growth potential with some risk
		}
	}
	
	// EPS and profitability
	if stock.EPS != nil && *stock.EPS > 0 {
		score += 10 // Positive earnings
	}
	
	// Dividend yield
	if stock.DividendYield != nil && *stock.DividendYield > 0 {
		yield := *stock.DividendYield
		if yield > 0.02 && yield < 0.08 { // 2-8% is attractive
			score += 5
		}
	}
	
	return score
}

func (s *RecommendationService) calculateTargetScore(stock domain.Stock) float64 {
	score := 0.0
	
	// Parse target prices and calculate upside
	targetFrom := s.parsePrice(stock.TargetFrom)
	targetTo := s.parsePrice(stock.TargetTo)
	
	var targetPrice float64
	if targetTo > 0 {
		targetPrice = targetTo
	} else if targetFrom > 0 {
		targetPrice = targetFrom
	}
	
	if targetPrice > 0 && stock.CurrentPrice != nil && *stock.CurrentPrice > 0 {
		currentPrice := *stock.CurrentPrice
		upside := (targetPrice - currentPrice) / currentPrice
		
		if upside > 0.3 { // >30% upside
			score += 25
		} else if upside > 0.15 { // 15-30% upside
			score += 15
		} else if upside > 0.05 { // 5-15% upside
			score += 8
		} else if upside < -0.1 { // >10% downside
			score -= 15
		}
	}
	
	return score
}

func (s *RecommendationService) parsePrice(priceStr string) float64 {
	if priceStr == "" {
		return 0
	}
	
	// Remove $ sign and any other non-numeric characters except decimal
	cleaned := strings.ReplaceAll(priceStr, "$", "")
	cleaned = strings.ReplaceAll(cleaned, ",", "")
	cleaned = strings.TrimSpace(cleaned)
	
	if val, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return val
	}
	
	return 0
}

func (s *RecommendationService) generateBasicReason(stock domain.Stock, score float64) string {
	reasons := []string{}
	
	// Add analyst action as primary reason
	if stock.Action != "" {
		reasons = append(reasons, fmt.Sprintf("analyst %s recommendation", strings.ToLower(stock.Action)))
	}
	
	// Add brokerage information
	if stock.Brokerage != "" {
		reasons = append(reasons, fmt.Sprintf("from %s", stock.Brokerage))
	}
	
	// Add rating change information
	if stock.RatingFrom != "" && stock.RatingTo != "" && stock.RatingFrom != stock.RatingTo {
		reasons = append(reasons, fmt.Sprintf("rating changed from %s to %s", stock.RatingFrom, stock.RatingTo))
	}
	
	// Add target price information if available
	if stock.TargetFrom != "" && stock.TargetTo != "" {
		reasons = append(reasons, fmt.Sprintf("price target: $%s to $%s", stock.TargetFrom, stock.TargetTo))
	}
	
	if len(reasons) == 0 {
		reasons = append(reasons, "based on current market analysis")
	}
	
	if score >= 70 {
		return fmt.Sprintf("Strong buy candidate %s. Score: %.1f/100", strings.Join(reasons, " "), score)
	} else if score >= 50 {
		return fmt.Sprintf("Moderate opportunity %s. Score: %.1f/100", strings.Join(reasons, " "), score)
	} else {
		return fmt.Sprintf("Hold consideration %s. Score: %.1f/100", strings.Join(reasons, " "), score)
	}
}

func (s *RecommendationService) getRecommendationType(score float64) string {
	if score >= 70 {
		return "BUY"
	} else if score >= 50 {
		return "HOLD"
	} else {
		return "WATCH"
	}
}