import { describe, it, expect, vi, beforeEach } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useStocksStore } from '../stocks'
import { apiService } from '@/services/api'

// Mock the API service
vi.mock('@/services/api', () => ({
  apiService: {
    getStocks: vi.fn(),
    syncStocks: vi.fn(),
    enrichStocks: vi.fn(),
    getRecommendations: vi.fn(),
    checkHealth: vi.fn()
  }
}))

describe('stocks store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.clearAllMocks()
  })

  it('initializes with default values', () => {
    const store = useStocksStore()
    
    expect(store.stocks).toEqual([])
    expect(store.recommendations).toEqual([])
    expect(store.enrichmentStats).toBeNull()
    expect(store.error).toBeNull()
    expect(store.loading.stocks).toBe(false)
    expect(store.loading.enrichment).toBe(false)
    expect(store.loading.recommendations).toBe(false)
    expect(store.loading.health).toBe(false)
    expect(store.loading.sync).toBe(false)
  })

  it('computes stocksCount correctly', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.' },
      { id: 2, ticker: 'GOOGL', company: 'Google LLC' }
    ] as any

    expect(store.stocksCount).toBe(2)
  })

  it('computes enrichedStocksCount correctly', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.', enriched_at: '2024-01-01' },
      { id: 2, ticker: 'GOOGL', company: 'Google LLC', enriched_at: null },
      { id: 3, ticker: 'MSFT', company: 'Microsoft', enriched_at: '2024-01-02' }
    ] as any

    expect(store.enrichedStocksCount).toBe(2)
  })

  it('computes sectorBreakdown correctly', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.', sector: 'Technology' },
      { id: 2, ticker: 'GOOGL', company: 'Google LLC', sector: 'Technology' },
      { id: 3, ticker: 'JPM', company: 'JPMorgan', sector: 'Financial Services' },
      { id: 4, ticker: 'TSLA', company: 'Tesla', sector: 'Automotive' }
    ] as any

    const breakdown = store.sectorBreakdown
    expect(breakdown).toHaveLength(3)
    expect(breakdown.find(s => s.sector === 'Technology')?.count).toBe(2)
    expect(breakdown.find(s => s.sector === 'Financial Services')?.count).toBe(1)
    expect(breakdown.find(s => s.sector === 'Automotive')?.count).toBe(1)
  })

  it('computes filteredStocks correctly with search filter', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.' },
      { id: 2, ticker: 'GOOGL', company: 'Google LLC' },
      { id: 3, ticker: 'MSFT', company: 'Microsoft Corp.' }
    ] as any

    store.updateFilters({ search: 'app' })
    
    expect(store.filteredStocks).toHaveLength(1)
    expect(store.filteredStocks[0].ticker).toBe('AAPL')
  })

  it('computes filteredStocks correctly with sector filter', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.', sector: 'Technology' },
      { id: 2, ticker: 'JPM', company: 'JPMorgan', sector: 'Financial Services' }
    ] as any

    store.updateFilters({ sector: 'Technology' })
    
    expect(store.filteredStocks).toHaveLength(1)
    expect(store.filteredStocks[0].ticker).toBe('AAPL')
  })

  it('computes filteredStocks correctly with enriched_only filter', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.', enriched_at: '2024-01-01' },
      { id: 2, ticker: 'GOOGL', company: 'Google LLC', enriched_at: null }
    ] as any

    store.updateFilters({ enriched_only: true })
    
    expect(store.filteredStocks).toHaveLength(1)
    expect(store.filteredStocks[0].ticker).toBe('AAPL')
  })

  it('computes paginatedStocks correctly', () => {
    const store = useStocksStore()
    
    store.stocks = Array.from({ length: 25 }, (_, i) => ({
      id: i + 1,
      ticker: `STOCK${i}`,
      company: `Company ${i}`
    })) as any

    // Default pagination (page 1, limit 20)
    expect(store.paginatedStocks).toHaveLength(20)
    expect(store.paginatedStocks[0].ticker).toBe('STOCK0')

    // Second page
    store.updatePagination({ page: 2 })
    expect(store.paginatedStocks).toHaveLength(5)
    expect(store.paginatedStocks[0].ticker).toBe('STOCK20')
  })

  it('fetchStocks sets loading state and handles success', async () => {
    const store = useStocksStore()
    const mockResponse = {
      stocks: [{ id: 1, ticker: 'AAPL', company: 'Apple Inc.' }],
      count: 1
    }

    vi.mocked(apiService.getStocks).mockResolvedValue(mockResponse)

    const promise = store.fetchStocks()
    expect(store.loading.stocks).toBe(true)

    await promise

    expect(store.loading.stocks).toBe(false)
    expect(store.stocks).toEqual(mockResponse.stocks)
    expect(store.pagination.total).toBe(1)
    expect(store.error).toBeNull()
  })

  it('fetchStocks handles API error', async () => {
    const store = useStocksStore()
    const errorMessage = 'API Error'

    vi.mocked(apiService.getStocks).mockRejectedValue(new Error(errorMessage))

    await store.fetchStocks()

    expect(store.loading.stocks).toBe(false)
    expect(store.error).toBe(errorMessage)
    expect(store.stocks).toEqual([])
  })

  it('fetchStocks skips API call if stocks exist and force is false', async () => {
    const store = useStocksStore()
    store.stocks = [{ id: 1, ticker: 'AAPL', company: 'Apple Inc.' }] as any

    await store.fetchStocks(false)

    expect(apiService.getStocks).not.toHaveBeenCalled()
  })

  it('syncStocks calls API and refreshes stocks', async () => {
    const store = useStocksStore()
    const mockStocksResponse = {
      stocks: [{ id: 1, ticker: 'AAPL', company: 'Apple Inc.' }],
      count: 1
    }

    vi.mocked(apiService.syncStocks).mockResolvedValue({ message: 'Success' })
    vi.mocked(apiService.getStocks).mockResolvedValue(mockStocksResponse)

    const promise = store.syncStocks()
    expect(store.loading.sync).toBe(true)

    await promise

    expect(store.loading.sync).toBe(false)
    expect(apiService.syncStocks).toHaveBeenCalled()
    expect(apiService.getStocks).toHaveBeenCalled()
    expect(store.error).toBeNull()
  })

  it('enrichStocks calls API and updates stats', async () => {
    const store = useStocksStore()
    const mockResponse = {
      message: 'Success',
      stats: { total_stocks: 5, newly_enriched: 3, already_enriched: 1, failed: 1 }
    }
    const mockStocksResponse = {
      stocks: [{ id: 1, ticker: 'AAPL', company: 'Apple Inc.' }],
      count: 1
    }

    vi.mocked(apiService.enrichStocks).mockResolvedValue(mockResponse)
    vi.mocked(apiService.getStocks).mockResolvedValue(mockStocksResponse)

    const promise = store.enrichStocks(5)
    expect(store.loading.enrichment).toBe(true)

    await promise

    expect(store.loading.enrichment).toBe(false)
    expect(apiService.enrichStocks).toHaveBeenCalledWith(5)
    expect(store.enrichmentStats).toEqual(mockResponse.stats)
    expect(store.error).toBeNull()
  })

  it('fetchRecommendations calls API and updates recommendations', async () => {
    const store = useStocksStore()
    const mockResponse = {
      recommendations: [
        { symbol: 'AAPL', name: 'Apple Inc.', score: 85, recommendation_type: 'BUY' }
      ],
      summary: { total_analyzed: 1, buy_recommendations: 1, hold_recommendations: 0 }
    }

    vi.mocked(apiService.getRecommendations).mockResolvedValue(mockResponse)

    const promise = store.fetchRecommendations(10)
    expect(store.loading.recommendations).toBe(true)

    await promise

    expect(store.loading.recommendations).toBe(false)
    expect(apiService.getRecommendations).toHaveBeenCalledWith(10)
    expect(store.recommendations).toEqual(mockResponse.recommendations)
    expect(store.error).toBeNull()
  })

  it('updateFilters updates filters and resets pagination', () => {
    const store = useStocksStore()
    
    store.pagination.page = 3
    store.updateFilters({ search: 'test', sector: 'Technology' })

    expect(store.filters.search).toBe('test')
    expect(store.filters.sector).toBe('Technology')
    expect(store.pagination.page).toBe(1)
  })

  it('updatePagination updates pagination correctly', () => {
    const store = useStocksStore()
    
    store.updatePagination({ page: 2, limit: 50 })

    expect(store.pagination.page).toBe(2)
    expect(store.pagination.limit).toBe(50)
  })

  it('clearFilters resets all filters and pagination', () => {
    const store = useStocksStore()
    
    store.filters.search = 'test'
    store.filters.sector = 'Technology'
    store.pagination.page = 3

    store.clearFilters()

    expect(store.filters.search).toBe('')
    expect(store.filters.sector).toBe('')
    expect(store.filters.recommendation_type).toBe('')
    expect(store.filters.action).toBe('')
    expect(store.filters.brokerage).toBe('')
    expect(store.filters.enriched_only).toBe(false)
    expect(store.pagination.page).toBe(1)
  })

  it('clearError clears the error state', () => {
    const store = useStocksStore()
    
    store.error = 'Test error'
    store.clearError()

    expect(store.error).toBeNull()
  })

  it('handles multiple filters simultaneously', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { 
        id: 1, 
        ticker: 'AAPL', 
        company: 'Apple Inc.', 
        sector: 'Technology',
        action: 'BUY',
        brokerage: 'Goldman Sachs',
        enriched_at: '2024-01-01'
      },
      { 
        id: 2, 
        ticker: 'GOOGL', 
        company: 'Google LLC', 
        sector: 'Technology',
        action: 'HOLD',
        brokerage: 'Morgan Stanley',
        enriched_at: null
      },
      { 
        id: 3, 
        ticker: 'JPM', 
        company: 'JPMorgan', 
        sector: 'Financial Services',
        action: 'BUY',
        brokerage: 'Goldman Sachs',
        enriched_at: '2024-01-01'
      }
    ] as any

    store.updateFilters({ 
      sector: 'Technology', 
      action: 'BUY',
      enriched_only: true
    })

    expect(store.filteredStocks).toHaveLength(1)
    expect(store.filteredStocks[0].ticker).toBe('AAPL')
  })

  it('handles case-insensitive search', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.' },
      { id: 2, ticker: 'googl', company: 'Google LLC' }
    ] as any

    store.updateFilters({ search: 'GOOGLE' })
    
    expect(store.filteredStocks).toHaveLength(1)
    expect(store.filteredStocks[0].ticker).toBe('googl')
  })

  it('handles empty sector breakdown when no stocks have sectors', () => {
    const store = useStocksStore()
    
    store.stocks = [
      { id: 1, ticker: 'AAPL', company: 'Apple Inc.' }, // no sector
      { id: 2, ticker: 'GOOGL', company: 'Google LLC' } // no sector
    ] as any

    expect(store.sectorBreakdown).toEqual([])
  })
})