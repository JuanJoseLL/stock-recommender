-- Drop indexes first
DROP INDEX IF EXISTS idx_stocks_ticker_time;
DROP INDEX IF EXISTS idx_stocks_brokerage;
DROP INDEX IF EXISTS idx_stocks_action;
DROP INDEX IF EXISTS idx_stocks_time;
DROP INDEX IF EXISTS idx_stocks_ticker;

-- Drop the table
DROP TABLE IF EXISTS stocks;