package repository

import (
    "context"
    "github.com/JuanJoseLL/stock-recommender/internal/domain"
)

// StockRepository defines stock data operations
type StockRepository interface {
    GetAll(ctx context.Context) ([]domain.Stock, error)
    GetByTicker(ctx context.Context, ticker string) (*domain.Stock, error)
    Create(ctx context.Context, stock *domain.Stock) error
    Update(ctx context.Context, stock *domain.Stock) error
    BulkCreate(ctx context.Context, stocks []domain.Stock) error
    GetTopRecommendations(ctx context.Context, limit int) ([]domain.Stock, error)
    UpsertStock(ctx context.Context, stock *domain.Stock) error
    CheckDuplicates(ctx context.Context) (int, error)
}