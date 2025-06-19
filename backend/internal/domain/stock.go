package domain

import "time"

type Stock struct {
    ID            int       `json:"id,omitempty" db:"id"`
    Ticker        string    `json:"ticker" db:"ticker"`
    TargetFrom    string    `json:"target_from" db:"target_from"`
    TargetTo      string    `json:"target_to" db:"target_to"`
    Company       string    `json:"company" db:"company"`
    Action        string    `json:"action" db:"action"`
    Brokerage     string    `json:"brokerage" db:"brokerage"`
    RatingFrom    string    `json:"rating_from" db:"rating_from"`
    RatingTo      string    `json:"rating_to" db:"rating_to"`
    Time          time.Time `json:"time" db:"time"`
    CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
    UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
    
    // Fundamental data from Alpha Vantage
    CurrentPrice   *float64 `json:"current_price,omitempty" db:"current_price"`
    Volume         *int64   `json:"volume,omitempty" db:"volume"`
    MarketCap      *int64   `json:"market_cap,omitempty" db:"market_cap"`
    PERatio        *float64 `json:"pe_ratio,omitempty" db:"pe_ratio"`
    WeekHigh52     *float64 `json:"week_high_52,omitempty" db:"week_high_52"`
    WeekLow52      *float64 `json:"week_low_52,omitempty" db:"week_low_52"`
    EPS            *float64 `json:"eps,omitempty" db:"eps"`
    BookValue      *float64 `json:"book_value,omitempty" db:"book_value"`
    DividendYield  *float64 `json:"dividend_yield,omitempty" db:"dividend_yield"`
    Sector         *string  `json:"sector,omitempty" db:"sector"`
    Industry       *string  `json:"industry,omitempty" db:"industry"`
    EnrichedAt     *time.Time `json:"enriched_at,omitempty" db:"enriched_at"`
}

type APIResponse struct {
    Items    []Stock `json:"items"`
    NextPage string  `json:"next_page,omitempty"`
}

type RecommendationResponse struct {
    RecommendedStock Stock   `json:"recommended_stock"`
    Confidence       float64 `json:"confidence"`
    Reason           string  `json:"reason"`
    AnalysisDate     string  `json:"analysis_date"`
}