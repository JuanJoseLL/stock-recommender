package service

import (
    "context"
    "fmt"
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

// SyncStocks fetches and saves stocks
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