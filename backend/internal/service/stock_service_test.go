package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JuanJoseLL/stock-recommender/internal/domain"
)

// Mock repository for testing
type mockStockRepository struct {
	stocks []domain.Stock
	err    error
}

func (m *mockStockRepository) GetAll(ctx context.Context) ([]domain.Stock, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stocks, nil
}

func (m *mockStockRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Stock, error) {
	return nil, nil
}

func (m *mockStockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	return m.err
}

func (m *mockStockRepository) Update(ctx context.Context, stock *domain.Stock) error {
	return m.err
}

func (m *mockStockRepository) BulkCreate(ctx context.Context, stocks []domain.Stock) error {
	return m.err
}

func (m *mockStockRepository) GetTopRecommendations(ctx context.Context, limit int) ([]domain.Stock, error) {
	return m.stocks, m.err
}

func (m *mockStockRepository) UpsertStock(ctx context.Context, stock *domain.Stock) error {
	return m.err
}

func (m *mockStockRepository) CheckDuplicates(ctx context.Context) (int, error) {
	return 0, m.err
}

// Mock external API client for testing
type mockExternalAPIClient struct {
	stocks []domain.Stock
	err    error
}

func (m *mockExternalAPIClient) FetchStocks(ctx context.Context) (*domain.APIResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.APIResponse{Items: m.stocks}, nil
}

func (m *mockExternalAPIClient) FetchStocksWithPagination(ctx context.Context, nextPage string) (*domain.APIResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &domain.APIResponse{Items: m.stocks}, nil
}

func (m *mockExternalAPIClient) FetchAllStocks(ctx context.Context) ([]domain.Stock, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stocks, nil
}

func TestNewStockService(t *testing.T) {
	repo := &mockStockRepository{}
	client := &mockExternalAPIClient{}

	service := NewStockService(repo, client)

	if service == nil {
		t.Error("Expected service to be created, got nil")
		return
	}
	if service.repo != repo {
		t.Error("Expected repository to be set correctly")
	}
	if service.client != client {
		t.Error("Expected client to be set correctly")
	}
}

func TestStockService_GetAllStocks(t *testing.T) {
	testStocks := []domain.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple Inc."},
		{ID: 2, Ticker: "GOOGL", Company: "Google LLC"},
	}

	repo := &mockStockRepository{stocks: testStocks}
	client := &mockExternalAPIClient{}
	service := NewStockService(repo, client)

	stocks, err := service.GetAllStocks(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(stocks) != 2 {
		t.Errorf("Expected 2 stocks, got %d", len(stocks))
	}
	if stocks[0].Ticker != "AAPL" {
		t.Errorf("Expected first stock ticker to be AAPL, got %s", stocks[0].Ticker)
	}
}

func TestStockService_GetAllStocks_Error(t *testing.T) {
	repo := &mockStockRepository{err: errors.New("database error")}
	client := &mockExternalAPIClient{}
	service := NewStockService(repo, client)

	stocks, err := service.GetAllStocks(context.Background())

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if stocks != nil {
		t.Error("Expected nil stocks on error")
	}
}

func TestStockService_GetStockStats(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)

	testStocks := []domain.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple Inc.", Brokerage: "Goldman Sachs", Time: now},
		{ID: 2, Ticker: "GOOGL", Company: "Google LLC", Brokerage: "Morgan Stanley", Time: yesterday},
		{ID: 3, Ticker: "AAPL", Company: "Apple Inc.", Brokerage: "Goldman Sachs", Time: now.Add(-12 * time.Hour)},
	}

	repo := &mockStockRepository{stocks: testStocks}
	client := &mockExternalAPIClient{}
	service := NewStockService(repo, client)

	stats, err := service.GetStockStats(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stats == nil {
		t.Error("Expected stats, got nil")
		return
	}

	if stats["total_stocks"] != 3 {
		t.Errorf("Expected total_stocks to be 3, got %v", stats["total_stocks"])
	}
	if stats["unique_tickers"] != 2 {
		t.Errorf("Expected unique_tickers to be 2, got %v", stats["unique_tickers"])
	}
	if stats["unique_companies"] != 2 {
		t.Errorf("Expected unique_companies to be 2, got %v", stats["unique_companies"])
	}
	if stats["unique_brokerages"] != 2 {
		t.Errorf("Expected unique_brokerages to be 2, got %v", stats["unique_brokerages"])
	}

	latestEntry, ok := stats["latest_entry"].(time.Time)
	if !ok {
		t.Error("Expected latest_entry to be time.Time")
	}
	if !latestEntry.Equal(now) {
		t.Errorf("Expected latest_entry to be %v, got %v", now, latestEntry)
	}

	oldestEntry, ok := stats["oldest_entry"].(time.Time)
	if !ok {
		t.Error("Expected oldest_entry to be time.Time")
	}
	if !oldestEntry.Equal(yesterday) {
		t.Errorf("Expected oldest_entry to be %v, got %v", yesterday, oldestEntry)
	}
}

func TestStockService_GetStockStats_EmptyStocks(t *testing.T) {
	repo := &mockStockRepository{stocks: []domain.Stock{}}
	client := &mockExternalAPIClient{}
	service := NewStockService(repo, client)

	stats, err := service.GetStockStats(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if stats["total_stocks"] != 0 {
		t.Errorf("Expected total_stocks to be 0, got %v", stats["total_stocks"])
	}
	if _, exists := stats["latest_entry"]; exists {
		t.Error("Expected latest_entry to not exist for empty stocks")
	}
}

func TestStockService_SyncStocks(t *testing.T) {
	testStocks := []domain.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple Inc."},
	}

	repo := &mockStockRepository{}
	client := &mockExternalAPIClient{stocks: testStocks}
	service := NewStockService(repo, client)

	err := service.SyncStocks(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestStockService_SyncStocks_FetchError(t *testing.T) {
	repo := &mockStockRepository{}
	client := &mockExternalAPIClient{err: errors.New("API error")}
	service := NewStockService(repo, client)

	err := service.SyncStocks(context.Background())

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if !errors.Is(err, errors.New("API error")) && err.Error() != "failed to fetch stocks: API error" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestStockService_SyncStocks_SaveError(t *testing.T) {
	testStocks := []domain.Stock{
		{ID: 1, Ticker: "AAPL", Company: "Apple Inc."},
	}

	repo := &mockStockRepository{err: errors.New("database error")}
	client := &mockExternalAPIClient{stocks: testStocks}
	service := NewStockService(repo, client)

	err := service.SyncStocks(context.Background())

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if err.Error() != "failed to save stocks: database error" {
		t.Errorf("Expected specific error message, got %v", err)
	}
}

func TestStockService_isValidStock(t *testing.T) {
	service := &StockService{}

	// Valid stock
	validStock := domain.Stock{
		Ticker:  "AAPL",
		Company: "Apple Inc.",
		Time:    time.Now(),
	}
	if !service.isValidStock(validStock) {
		t.Error("Expected valid stock to pass validation")
	}

	// Invalid stock - empty ticker
	invalidStock1 := domain.Stock{
		Company: "Apple Inc.",
		Time:    time.Now(),
	}
	if service.isValidStock(invalidStock1) {
		t.Error("Expected stock with empty ticker to fail validation")
	}

	// Invalid stock - empty company
	invalidStock2 := domain.Stock{
		Ticker: "AAPL",
		Time:   time.Now(),
	}
	if service.isValidStock(invalidStock2) {
		t.Error("Expected stock with empty company to fail validation")
	}

	// Invalid stock - zero time
	invalidStock3 := domain.Stock{
		Ticker:  "AAPL",
		Company: "Apple Inc.",
	}
	if service.isValidStock(invalidStock3) {
		t.Error("Expected stock with zero time to fail validation")
	}
}

func TestStockService_validateStocks(t *testing.T) {
	service := &StockService{}

	stocks := []domain.Stock{
		{Ticker: "AAPL", Company: "Apple Inc.", Time: time.Now()}, // Valid
		{Ticker: "", Company: "Google", Time: time.Now()},         // Invalid
		{Ticker: "MSFT", Company: "Microsoft", Time: time.Now()},  // Valid
		{Ticker: "TSLA", Company: "", Time: time.Now()},           // Invalid
		{Ticker: "AMZN", Company: "Amazon", Time: time.Time{}},    // Invalid
	}

	validStocks := service.validateStocks(stocks)

	if len(validStocks) != 2 {
		t.Errorf("Expected 2 valid stocks, got %d", len(validStocks))
	}
	if validStocks[0].Ticker != "AAPL" {
		t.Errorf("Expected first valid stock to be AAPL, got %s", validStocks[0].Ticker)
	}
	if validStocks[1].Ticker != "MSFT" {
		t.Errorf("Expected second valid stock to be MSFT, got %s", validStocks[1].Ticker)
	}
}

func TestStockService_countUniqueTickers(t *testing.T) {
	service := &StockService{}

	stocks := []domain.Stock{
		{Ticker: "AAPL"},
		{Ticker: "GOOGL"},
		{Ticker: "AAPL"}, // Duplicate
		{Ticker: "MSFT"},
	}

	count := service.countUniqueTickers(stocks)
	if count != 3 {
		t.Errorf("Expected 3 unique tickers, got %d", count)
	}
}

func TestStockService_countUniqueCompanies(t *testing.T) {
	service := &StockService{}

	stocks := []domain.Stock{
		{Company: "Apple Inc."},
		{Company: "Google LLC"},
		{Company: "Apple Inc."}, // Duplicate
		{Company: "Microsoft"},
	}

	count := service.countUniqueCompanies(stocks)
	if count != 3 {
		t.Errorf("Expected 3 unique companies, got %d", count)
	}
}

func TestStockService_countUniqueBrokerages(t *testing.T) {
	service := &StockService{}

	stocks := []domain.Stock{
		{Brokerage: "Goldman Sachs"},
		{Brokerage: "Morgan Stanley"},
		{Brokerage: "Goldman Sachs"}, // Duplicate
		{Brokerage: "JP Morgan"},
	}

	count := service.countUniqueBrokerages(stocks)
	if count != 3 {
		t.Errorf("Expected 3 unique brokerages, got %d", count)
	}
}

func TestStockService_getLatestEntry(t *testing.T) {
	service := &StockService{}

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	tomorrow := now.Add(24 * time.Hour)

	stocks := []domain.Stock{
		{Time: yesterday},
		{Time: now},
		{Time: tomorrow},
	}

	latest := service.getLatestEntry(stocks)
	if !latest.Equal(tomorrow) {
		t.Errorf("Expected latest entry to be %v, got %v", tomorrow, latest)
	}
}

func TestStockService_getOldestEntry(t *testing.T) {
	service := &StockService{}

	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	dayBeforeYesterday := now.Add(-48 * time.Hour)

	stocks := []domain.Stock{
		{Time: yesterday},
		{Time: now},
		{Time: dayBeforeYesterday},
	}

	oldest := service.getOldestEntry(stocks)
	if !oldest.Equal(dayBeforeYesterday) {
		t.Errorf("Expected oldest entry to be %v, got %v", dayBeforeYesterday, oldest)
	}
}
