package service

import (
    "context"
    "fmt"
    "log"
    "time"
    "github.com/JuanJoseLL/stock-recommender/internal/client"
    "github.com/JuanJoseLL/stock-recommender/internal/domain"
    "github.com/JuanJoseLL/stock-recommender/internal/repository"
)

// StockService handles business logic
type StockService struct {
    repo   repository.StockRepository
    client client.ExternalAPIClient
}

// NewStockService creates a new stock service
func NewStockService(repo repository.StockRepository, client client.ExternalAPIClient) *StockService {
    return &StockService{
        repo:   repo,
        client: client,
    }
}

// GetAllStocks retrieves all stocks
func (s *StockService) GetAllStocks(ctx context.Context) ([]domain.Stock, error) {
    return s.repo.GetAll(ctx)
}

// GetStockStats returns statistics about stored stocks
func (s *StockService) GetStockStats(ctx context.Context) (map[string]interface{}, error) {
    allStocks, err := s.repo.GetAll(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to get stocks for stats: %w", err)
    }
    
    stats := map[string]interface{}{
        "total_stocks": len(allStocks),
        "unique_tickers": s.countUniqueTickers(allStocks),
        "unique_companies": s.countUniqueCompanies(allStocks),
        "unique_brokerages": s.countUniqueBrokerages(allStocks),
    }
    
    if len(allStocks) > 0 {
        stats["latest_entry"] = s.getLatestEntry(allStocks)
        stats["oldest_entry"] = s.getOldestEntry(allStocks)
    }
    
    return stats, nil
}

func (s *StockService) countUniqueTickers(stocks []domain.Stock) int {
    tickers := make(map[string]bool)
    for _, stock := range stocks {
        tickers[stock.Ticker] = true
    }
    return len(tickers)
}

func (s *StockService) countUniqueCompanies(stocks []domain.Stock) int {
    companies := make(map[string]bool)
    for _, stock := range stocks {
        companies[stock.Company] = true
    }
    return len(companies)
}

func (s *StockService) countUniqueBrokerages(stocks []domain.Stock) int {
    brokerages := make(map[string]bool)
    for _, stock := range stocks {
        brokerages[stock.Brokerage] = true
    }
    return len(brokerages)
}

func (s *StockService) getLatestEntry(stocks []domain.Stock) time.Time {
    var latest time.Time
    for _, stock := range stocks {
        if stock.Time.After(latest) {
            latest = stock.Time
        }
    }
    return latest
}

func (s *StockService) getOldestEntry(stocks []domain.Stock) time.Time {
    oldest := time.Now()
    for _, stock := range stocks {
        if stock.Time.Before(oldest) {
            oldest = stock.Time
        }
    }
    return oldest
}

// SyncStocks fetches and saves stocks (single page - legacy method)
func (s *StockService) SyncStocks(ctx context.Context) error {
    apiResponse, err := s.client.FetchStocks(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch stocks: %w", err)
    }
    
    if err := s.repo.BulkCreate(ctx, apiResponse.Items); err != nil {
        return fmt.Errorf("failed to save stocks: %w", err)
    }
    
    return nil
}

// SyncAllStocks fetches and saves all stocks using pagination
func (s *StockService) SyncAllStocks(ctx context.Context) error {
    log.Println("Starting complete stock synchronization...")
    startTime := time.Now()
    
    allStocks, err := s.client.FetchAllStocks(ctx)
    if err != nil {
        return fmt.Errorf("failed to fetch all stocks: %w", err)
    }
    
    if len(allStocks) == 0 {
        log.Println("No stocks received from API")
        return nil
    }
    
    log.Printf("Fetched %d stocks from API", len(allStocks))
    
    // Validate stocks before saving
    validStocks := s.validateStocks(allStocks)
    if len(validStocks) != len(allStocks) {
        log.Printf("Filtered out %d invalid stocks", len(allStocks)-len(validStocks))
    }
    
    if len(validStocks) == 0 {
        return fmt.Errorf("no valid stocks to save")
    }
    
    // Save stocks in batches for better performance
    const batchSize = 100
    totalSaved := 0
    
    for i := 0; i < len(validStocks); i += batchSize {
        end := i + batchSize
        if end > len(validStocks) {
            end = len(validStocks)
        }
        
        batch := validStocks[i:end]
        if err := s.repo.BulkCreate(ctx, batch); err != nil {
            return fmt.Errorf("failed to save batch %d-%d: %w", i, end-1, err)
        }
        
        totalSaved += len(batch)
        log.Printf("Saved batch %d-%d (%d/%d stocks)", i, end-1, totalSaved, len(validStocks))
        
        // Small delay between batches to avoid overwhelming the database
        if end < len(validStocks) {
            time.Sleep(50 * time.Millisecond)
        }
    }
    
    duration := time.Since(startTime)
    log.Printf("Stock synchronization completed: %d stocks saved in %v", totalSaved, duration)
    
    return nil
}

// validateStocks filters out invalid stock entries
func (s *StockService) validateStocks(stocks []domain.Stock) []domain.Stock {
    var validStocks []domain.Stock
    
    for _, stock := range stocks {
        if s.isValidStock(stock) {
            validStocks = append(validStocks, stock)
        }
    }
    
    return validStocks
}

// isValidStock checks if a stock entry has required fields
func (s *StockService) isValidStock(stock domain.Stock) bool {
    if stock.Ticker == "" {
        return false
    }
    if stock.Company == "" {
        return false
    }
    if stock.Time.IsZero() {
        return false
    }
    return true
}