CREATE TABLE IF NOT EXISTS stocks (
    id SERIAL PRIMARY KEY,
    ticker VARCHAR(10) NOT NULL,
    target_from VARCHAR(20),
    target_to VARCHAR(20),
    company VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    brokerage VARCHAR(100) NOT NULL,
    rating_from VARCHAR(50),
    rating_to VARCHAR(50),
    time TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_stocks_ticker ON stocks(ticker);
CREATE INDEX IF NOT EXISTS idx_stocks_time ON stocks(time DESC);
CREATE INDEX IF NOT EXISTS idx_stocks_action ON stocks(action);
CREATE INDEX IF NOT EXISTS idx_stocks_brokerage ON stocks(brokerage);

-- Create composite index for common queries
CREATE INDEX IF NOT EXISTS idx_stocks_ticker_time ON stocks(ticker, time DESC);