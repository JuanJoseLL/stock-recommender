export interface Stock {
  id: number
  ticker: string
  target_from: string
  target_to: string
  company: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
  market_cap?: number
  pe_ratio?: number
  week_high_52?: number
  week_low_52?: number
  eps?: number
  sector?: string
  industry?: string
  enriched_at?: string
}

export interface StocksResponse {
  stocks: Stock[]
  count: number
}

export interface Recommendation {
  symbol: string
  name: string
  score: number
  reason: string
  current_price?: string
  target_price?: string
  pe_ratio?: string
  market_cap?: string
  sector?: string
  recommendation_type: 'BUY' | 'HOLD' | 'WATCH' | 'SELL'
}

export interface RecommendationsResponse {
  recommendations: Recommendation[]
  summary: {
    total_analyzed: number
    buy_recommendations: number
    hold_recommendations: number
    generated_at: string
    data_source: string
  }
}

export interface EnrichmentStats {
  total_stocks: number
  already_enriched: number
  newly_enriched: number
  failed: number
  rate_limit_reached: boolean
}

export interface EnrichmentResponse {
  message: string
  stats: EnrichmentStats
}

export interface HealthResponse {
  status: string
  service: string
}

export interface SyncResponse {
  message: string
}

export interface StockDetail extends Stock {
  upside_potential?: number
  risk_level?: 'Low' | 'Medium' | 'High'
  analyst_consensus?: string
}

export interface DashboardData {
  totalStocks: number
  enrichedStocks: number
  topRecommendations: Recommendation[]
  sectorBreakdown: { sector: string; count: number }[]
}

export interface AppFilters {
  search: string
  sector: string
  recommendation_type: string
  action: string
  brokerage: string
  enriched_only: boolean
}

export interface AppPagination {
  page: number
  limit: number
  total: number
}

export interface AppState {
  stocks: Stock[]
  recommendations: Recommendation[]
  loading: {
    stocks: boolean
    enrichment: boolean
    recommendations: boolean
    health: boolean
  }
  filters: AppFilters
  pagination: AppPagination
  error: string | null
}