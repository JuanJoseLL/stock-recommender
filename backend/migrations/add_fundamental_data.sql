-- Add fundamental data columns to stocks table
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS current_price DECIMAL(10,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS volume BIGINT;
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS market_cap BIGINT;
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS pe_ratio DECIMAL(8,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS week_high_52 DECIMAL(10,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS week_low_52 DECIMAL(10,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS eps DECIMAL(8,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS book_value DECIMAL(8,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS dividend_yield DECIMAL(6,4);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS sector VARCHAR(100);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS industry VARCHAR(150);
ALTER TABLE stocks ADD COLUMN IF NOT EXISTS enriched_at TIMESTAMP;

-- Create index on enriched_at for efficient querying of enriched stocks
CREATE INDEX IF NOT EXISTS idx_stocks_enriched_at ON stocks(enriched_at);

-- Create index on sector for filtering
CREATE INDEX IF NOT EXISTS idx_stocks_sector ON stocks(sector);

-- Create index on pe_ratio for sorting/filtering
CREATE INDEX IF NOT EXISTS idx_stocks_pe_ratio ON stocks(pe_ratio);