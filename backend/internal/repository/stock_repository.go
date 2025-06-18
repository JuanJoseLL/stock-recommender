package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/JuanJoseLL/stock-recommender/internal/domain"
	"github.com/JuanJoseLL/stock-recommender/pkg/database"
)

type stockRepository struct {
	db *database.DB
}

func NewStockRepository(db *database.DB) StockRepository {
	return &stockRepository{
		db: db,
	}
}

func (r *stockRepository) GetAll(ctx context.Context) ([]domain.Stock, error) {
	query := `
		SELECT id, ticker, target_from, target_to, company, action, 
		       brokerage, rating_from, rating_to, time, created_at, updated_at
		FROM stocks
		ORDER BY time DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stocks: %w", err)
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(
			&stock.ID,
			&stock.Ticker,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Company,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan stock row: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating stock rows: %w", err)
	}

	return stocks, nil
}

func (r *stockRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Stock, error) {
	query := `
		SELECT id, ticker, target_from, target_to, company, action, 
		       brokerage, rating_from, rating_to, time, created_at, updated_at
		FROM stocks
		WHERE ticker = $1
		ORDER BY time DESC
		LIMIT 1
	`
	
	var stock domain.Stock
	err := r.db.QueryRowContext(ctx, query, ticker).Scan(
		&stock.ID,
		&stock.Ticker,
		&stock.TargetFrom,
		&stock.TargetTo,
		&stock.Company,
		&stock.Action,
		&stock.Brokerage,
		&stock.RatingFrom,
		&stock.RatingTo,
		&stock.Time,
		&stock.CreatedAt,
		&stock.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Stock not found
		}
		return nil, fmt.Errorf("failed to get stock by ticker %s: %w", ticker, err)
	}

	return &stock, nil
}

func (r *stockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	query := `
		INSERT INTO stocks (ticker, target_from, target_to, company, action, 
		                   brokerage, rating_from, rating_to, time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRowContext(ctx, query,
		stock.Ticker,
		stock.TargetFrom,
		stock.TargetTo,
		stock.Company,
		stock.Action,
		stock.Brokerage,
		stock.RatingFrom,
		stock.RatingTo,
		stock.Time,
	).Scan(&stock.ID, &stock.CreatedAt, &stock.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create stock: %w", err)
	}

	log.Printf("Created stock with ID: %d, Ticker: %s", stock.ID, stock.Ticker)
	return nil
}

func (r *stockRepository) BulkCreate(ctx context.Context, stocks []domain.Stock) error {
	if len(stocks) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Use UPSERT to handle duplicates - insert or update on conflict
	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO stocks (ticker, target_from, target_to, company, action, 
		                   brokerage, rating_from, rating_to, time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (ticker, time) DO UPDATE SET
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			company = EXCLUDED.company,
			action = EXCLUDED.action,
			brokerage = EXCLUDED.brokerage,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			updated_at = NOW()
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare bulk upsert statement: %w", err)
	}
	defer stmt.Close()

	successCount := 0
	for _, stock := range stocks {
		_, err := stmt.ExecContext(ctx,
			stock.Ticker,
			stock.TargetFrom,
			stock.TargetTo,
			stock.Company,
			stock.Action,
			stock.Brokerage,
			stock.RatingFrom,
			stock.RatingTo,
			stock.Time,
		)
		if err != nil {
			// Log the error but continue with other stocks
			log.Printf("Failed to upsert stock %s: %v", stock.Ticker, err)
			continue
		}
		successCount++
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit bulk upsert transaction: %w", err)
	}

	log.Printf("Successfully upserted %d/%d stocks", successCount, len(stocks))
	return nil
}

func (r *stockRepository) GetTopRecommendations(ctx context.Context, limit int) ([]domain.Stock, error) {
	query := `
		SELECT id, ticker, target_from, target_to, company, action, 
		       brokerage, rating_from, rating_to, time, created_at, updated_at
		FROM stocks
		WHERE action IN ('Buy', 'Strong Buy', 'Upgrade')
		ORDER BY time DESC
		LIMIT $1
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top recommendations: %w", err)
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var stock domain.Stock
		err := rows.Scan(
			&stock.ID,
			&stock.Ticker,
			&stock.TargetFrom,
			&stock.TargetTo,
			&stock.Company,
			&stock.Action,
			&stock.Brokerage,
			&stock.RatingFrom,
			&stock.RatingTo,
			&stock.Time,
			&stock.CreatedAt,
			&stock.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recommendation row: %w", err)
		}
		stocks = append(stocks, stock)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recommendation rows: %w", err)
	}

	return stocks, nil
}

// UpsertStock inserts or updates a stock record
func (r *stockRepository) UpsertStock(ctx context.Context, stock *domain.Stock) error {
	query := `
		INSERT INTO stocks (ticker, target_from, target_to, company, action, 
		                   brokerage, rating_from, rating_to, time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		ON CONFLICT (ticker, time) DO UPDATE SET
			target_from = EXCLUDED.target_from,
			target_to = EXCLUDED.target_to,
			company = EXCLUDED.company,
			action = EXCLUDED.action,
			brokerage = EXCLUDED.brokerage,
			rating_from = EXCLUDED.rating_from,
			rating_to = EXCLUDED.rating_to,
			updated_at = NOW()
		RETURNING id, created_at, updated_at
	`
	
	err := r.db.QueryRowContext(ctx, query,
		stock.Ticker,
		stock.TargetFrom,
		stock.TargetTo,
		stock.Company,
		stock.Action,
		stock.Brokerage,
		stock.RatingFrom,
		stock.RatingTo,
		stock.Time,
	).Scan(&stock.ID, &stock.CreatedAt, &stock.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to upsert stock: %w", err)
	}

	return nil
}

// CheckDuplicates returns the count of potential duplicates
func (r *stockRepository) CheckDuplicates(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*) as duplicate_count
		FROM (
			SELECT ticker, time, COUNT(*) as cnt
			FROM stocks
			GROUP BY ticker, time
			HAVING COUNT(*) > 1
		) as duplicates
	`
	
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to check duplicates: %w", err)
	}
	
	return count, nil
}