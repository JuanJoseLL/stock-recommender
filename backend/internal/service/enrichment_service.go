package service

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/JuanJoseLL/stock-recommender/internal/client"
	"github.com/JuanJoseLL/stock-recommender/internal/domain"
	"github.com/JuanJoseLL/stock-recommender/internal/repository"
)

type EnrichmentService struct {
	stockRepo          repository.StockRepository
	alphaVantageClient *client.AlphaVantageClient
}

type EnrichmentStats struct {
	TotalStocks       int `json:"total_stocks"`
	AlreadyEnriched   int `json:"already_enriched"`
	NewlyEnriched     int `json:"newly_enriched"`
	Failed            int `json:"failed"`
	RateLimitReached  bool `json:"rate_limit_reached"`
}

func NewEnrichmentService(stockRepo repository.StockRepository, alphaVantageClient *client.AlphaVantageClient) *EnrichmentService {
	return &EnrichmentService{
		stockRepo:          stockRepo,
		alphaVantageClient: alphaVantageClient,
	}
}

func (s *EnrichmentService) EnrichStockData(ctx context.Context, limit int) (*EnrichmentStats, error) {
	stats := &EnrichmentStats{}
	
	// Get unique tickers that need enrichment
	tickers, err := s.getTickersForEnrichment(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get tickers for enrichment: %w", err)
	}
	
	stats.TotalStocks = len(tickers)
	log.Printf("Starting enrichment for %d unique tickers", len(tickers))
	
	for i, ticker := range tickers {
		select {
		case <-ctx.Done():
			return stats, ctx.Err()
		default:
		}
		
		log.Printf("Enriching ticker %s (%d/%d)", ticker, i+1, len(tickers))
		
		// Check if already enriched recently (within 24 hours)
		if s.isRecentlyEnriched(ctx, ticker) {
			stats.AlreadyEnriched++
			log.Printf("Ticker %s already enriched recently, skipping", ticker)
			continue
		}
		
		// Get company overview from Alpha Vantage
		overview, err := s.alphaVantageClient.GetCompanyOverview(ctx, ticker)
		if err != nil {
			log.Printf("Failed to get overview for %s: %v", ticker, err)
			stats.Failed++
			
			// Check if it's a rate limit error
			if strings.Contains(err.Error(), "Thank you for using Alpha Vantage") ||
				strings.Contains(err.Error(), "higher API call frequency") {
				stats.RateLimitReached = true
				log.Println("Rate limit reached, stopping enrichment")
				break
			}
			continue
		}
		
		// Convert overview to fundamental data
		fundamentals := s.convertOverviewToFundamentals(overview)
		
		// Update all stocks with this ticker
		err = s.updateStockFundamentals(ctx, ticker, fundamentals)
		if err != nil {
			log.Printf("Failed to update fundamentals for %s: %v", ticker, err)
			stats.Failed++
			continue
		}
		
		stats.NewlyEnriched++
		log.Printf("Successfully enriched ticker %s", ticker)
		
		// Rate limiting: wait between requests (Alpha Vantage allows 5 requests per minute)
		if i < len(tickers)-1 {
			time.Sleep(15 * time.Second) // 4 requests per minute to be safe
		}
	}
	
	log.Printf("Enrichment completed. Stats: %+v", stats)
	return stats, nil
}

func (s *EnrichmentService) getTickersForEnrichment(ctx context.Context, limit int) ([]string, error) {
	// Get all stocks from database
	stocks, err := s.stockRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	// Get unique tickers that haven't been enriched or are old
	tickerMap := make(map[string]bool)
	var tickers []string
	
	for _, stock := range stocks {
		if stock.Ticker == "" {
			continue
		}
		
		// Skip if already processed this ticker
		if tickerMap[stock.Ticker] {
			continue
		}
		
		// Add ticker if not enriched or enriched more than 24 hours ago
		if stock.EnrichedAt == nil || time.Since(*stock.EnrichedAt) > 24*time.Hour {
			tickers = append(tickers, stock.Ticker)
			tickerMap[stock.Ticker] = true
		}
		
		// Apply limit
		if len(tickers) >= limit {
			break
		}
	}
	
	return tickers, nil
}

func (s *EnrichmentService) isRecentlyEnriched(ctx context.Context, ticker string) bool {
	stocks, err := s.stockRepo.GetAll(ctx)
	if err != nil {
		return false
	}
	
	for _, stock := range stocks {
		if stock.Ticker == ticker && stock.EnrichedAt != nil {
			return time.Since(*stock.EnrichedAt) < 24*time.Hour
		}
	}
	
	return false
}

func (s *EnrichmentService) convertOverviewToFundamentals(overview *client.AlphaVantageCompanyOverview) map[string]interface{} {
	fundamentals := make(map[string]interface{})
	now := time.Now()
	
	// Parse and convert numeric fields
	if val := s.parseFloat(overview.MarketCapitalization); val != nil {
		fundamentals["market_cap"] = int64(*val)
	}
	
	if val := s.parseFloat(overview.PERatio); val != nil {
		fundamentals["pe_ratio"] = *val
	}
	
	if val := s.parseFloat(overview.WeekHigh52); val != nil {
		fundamentals["week_high_52"] = *val
	}
	
	if val := s.parseFloat(overview.WeekLow52); val != nil {
		fundamentals["week_low_52"] = *val
	}
	
	if val := s.parseFloat(overview.EPS); val != nil {
		fundamentals["eps"] = *val
	}
	
	if val := s.parseFloat(overview.BookValue); val != nil {
		fundamentals["book_value"] = *val
	}
	
	if val := s.parseFloat(overview.DividendYield); val != nil {
		fundamentals["dividend_yield"] = *val
	}
	
	// String fields
	if overview.Sector != "" && overview.Sector != "None" {
		fundamentals["sector"] = overview.Sector
	}
	
	if overview.Industry != "" && overview.Industry != "None" {
		fundamentals["industry"] = overview.Industry
	}
	
	// Set enrichment timestamp
	fundamentals["enriched_at"] = now
	
	return fundamentals
}

func (s *EnrichmentService) parseFloat(str string) *float64 {
	if str == "" || str == "None" || str == "-" {
		return nil
	}
	
	// Remove any non-numeric characters except decimal point and minus sign
	cleaned := strings.ReplaceAll(str, ",", "")
	cleaned = strings.ReplaceAll(cleaned, "$", "")
	cleaned = strings.TrimSpace(cleaned)
	
	if val, err := strconv.ParseFloat(cleaned, 64); err == nil {
		return &val
	}
	
	return nil
}

func (s *EnrichmentService) updateStockFundamentals(ctx context.Context, ticker string, fundamentals map[string]interface{}) error {
	// Get all stocks with this ticker
	allStocks, err := s.stockRepo.GetAll(ctx)
	if err != nil {
		return err
	}
	
	// Filter stocks by ticker
	var tickerStocks []domain.Stock
	for _, stock := range allStocks {
		if stock.Ticker == ticker {
			tickerStocks = append(tickerStocks, stock)
		}
	}
	
	if len(tickerStocks) == 0 {
		return fmt.Errorf("no stocks found for ticker %s", ticker)
	}
	
	// Update each stock record with fundamental data
	for _, stock := range tickerStocks {
		updatedStock := stock
		
		// Apply fundamental data
		if val, ok := fundamentals["market_cap"]; ok {
			if intVal, ok := val.(int64); ok {
				updatedStock.MarketCap = &intVal
			}
		}
		
		if val, ok := fundamentals["pe_ratio"]; ok {
			if floatVal, ok := val.(float64); ok {
				updatedStock.PERatio = &floatVal
			}
		}
		
		if val, ok := fundamentals["week_high_52"]; ok {
			if floatVal, ok := val.(float64); ok {
				updatedStock.WeekHigh52 = &floatVal
			}
		}
		
		if val, ok := fundamentals["week_low_52"]; ok {
			if floatVal, ok := val.(float64); ok {
				updatedStock.WeekLow52 = &floatVal
			}
		}
		
		if val, ok := fundamentals["eps"]; ok {
			if floatVal, ok := val.(float64); ok {
				updatedStock.EPS = &floatVal
			}
		}
		
		if val, ok := fundamentals["book_value"]; ok {
			if floatVal, ok := val.(float64); ok {
				updatedStock.BookValue = &floatVal
			}
		}
		
		if val, ok := fundamentals["dividend_yield"]; ok {
			if floatVal, ok := val.(float64); ok {
				updatedStock.DividendYield = &floatVal
			}
		}
		
		if val, ok := fundamentals["sector"]; ok {
			if strVal, ok := val.(string); ok {
				updatedStock.Sector = &strVal
			}
		}
		
		if val, ok := fundamentals["industry"]; ok {
			if strVal, ok := val.(string); ok {
				updatedStock.Industry = &strVal
			}
		}
		
		if val, ok := fundamentals["enriched_at"]; ok {
			if timeVal, ok := val.(time.Time); ok {
				updatedStock.EnrichedAt = &timeVal
			}
		}
		
		// Update in database
		if err := s.stockRepo.Update(ctx, &updatedStock); err != nil {
			return fmt.Errorf("failed to update stock %d: %w", stock.ID, err)
		}
	}
	
	return nil
}