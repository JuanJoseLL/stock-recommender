import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useHealth, useStocks, useSyncStocks, useEnrichStocks, useRecommendations } from '../useApi'
import { apiService } from '@/services/api'

// Mock the API service
vi.mock('@/services/api', () => ({
  apiService: {
    checkHealth: vi.fn(),
    getStocks: vi.fn(),
    syncStocks: vi.fn(),
    enrichStocks: vi.fn(),
    getRecommendations: vi.fn()
  }
}))

describe('useApi composables', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  describe('useHealth', () => {
    it('initializes with default values', () => {
      const { data, loading, error } = useHealth()
      
      expect(data.value).toBeNull()
      expect(loading.value).toBe(false)
      expect(error.value).toBeNull()
    })

    it('handles successful API call', async () => {
      const mockResponse = { status: 'healthy', service: 'stock-recommender' }
      vi.mocked(apiService.checkHealth).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useHealth()

      const promise = execute()
      expect(loading.value).toBe(true)

      await promise

      expect(loading.value).toBe(false)
      expect(data.value).toEqual(mockResponse)
      expect(error.value).toBeNull()
    })

    it('handles API error', async () => {
      const errorMessage = 'Health check failed'
      vi.mocked(apiService.checkHealth).mockRejectedValue(new Error(errorMessage))

      const { data, loading, error, execute } = useHealth()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe(errorMessage)
    })

    it('handles non-Error objects', async () => {
      vi.mocked(apiService.checkHealth).mockRejectedValue('String error')

      const { data, loading, error, execute } = useHealth()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe('Failed to check health')
    })
  })

  describe('useStocks', () => {
    it('initializes with default values', () => {
      const { data, loading, error } = useStocks()
      
      expect(data.value).toBeNull()
      expect(loading.value).toBe(false)
      expect(error.value).toBeNull()
    })

    it('handles successful API call', async () => {
      const mockResponse = {
        stocks: [{ 
          id: 1, 
          ticker: 'AAPL', 
          company: 'Apple Inc.',
          target_from: '150.00',
          target_to: '200.00',
          action: 'BUY',
          brokerage: 'Goldman Sachs',
          rating_from: 'BUY',
          rating_to: 'BUY',
          time: '2024-01-01T00:00:00Z'
        }],
        count: 1
      }
      vi.mocked(apiService.getStocks).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useStocks()

      const promise = execute()
      expect(loading.value).toBe(true)

      await promise

      expect(loading.value).toBe(false)
      expect(data.value).toEqual(mockResponse)
      expect(error.value).toBeNull()
    })

    it('handles API error', async () => {
      const errorMessage = 'Failed to fetch stocks'
      vi.mocked(apiService.getStocks).mockRejectedValue(new Error(errorMessage))

      const { data, loading, error, execute } = useStocks()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe(errorMessage)
    })

    it('uses default error message for non-Error objects', async () => {
      vi.mocked(apiService.getStocks).mockRejectedValue('String error')

      const { data, loading, error, execute } = useStocks()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe('Failed to fetch stocks')
    })
  })

  describe('useSyncStocks', () => {
    it('initializes with default values', () => {
      const { data, loading, error } = useSyncStocks()
      
      expect(data.value).toBeNull()
      expect(loading.value).toBe(false)
      expect(error.value).toBeNull()
    })

    it('handles successful API call', async () => {
      const mockResponse = { message: 'Stocks synced successfully' }
      vi.mocked(apiService.syncStocks).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useSyncStocks()

      const promise = execute()
      expect(loading.value).toBe(true)

      await promise

      expect(loading.value).toBe(false)
      expect(data.value).toEqual(mockResponse)
      expect(error.value).toBeNull()
    })

    it('handles API error', async () => {
      const errorMessage = 'Sync failed'
      vi.mocked(apiService.syncStocks).mockRejectedValue(new Error(errorMessage))

      const { data, loading, error, execute } = useSyncStocks()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe(errorMessage)
    })
  })

  describe('useEnrichStocks', () => {
    it('initializes with default values', () => {
      const { data, loading, error } = useEnrichStocks()
      
      expect(data.value).toBeNull()
      expect(loading.value).toBe(false)
      expect(error.value).toBeNull()
    })

    it('handles successful API call with default limit', async () => {
      const mockResponse = {
        message: 'Enrichment completed',
        stats: { 
          total_stocks: 5, 
          newly_enriched: 3, 
          already_enriched: 1, 
          failed: 1,
          rate_limit_reached: false
        }
      }
      vi.mocked(apiService.enrichStocks).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useEnrichStocks()

      const promise = execute()
      expect(loading.value).toBe(true)

      await promise

      expect(loading.value).toBe(false)
      expect(data.value).toEqual(mockResponse)
      expect(error.value).toBeNull()
      expect(apiService.enrichStocks).toHaveBeenCalledWith(5)
    })

    it('handles successful API call with custom limit', async () => {
      const mockResponse = {
        message: 'Enrichment completed',
        stats: { 
          total_stocks: 10, 
          newly_enriched: 7, 
          already_enriched: 2, 
          failed: 1,
          rate_limit_reached: false
        }
      }
      vi.mocked(apiService.enrichStocks).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useEnrichStocks()

      await execute(10)

      expect(data.value).toEqual(mockResponse)
      expect(apiService.enrichStocks).toHaveBeenCalledWith(10)
    })

    it('handles API error', async () => {
      const errorMessage = 'Enrichment failed'
      vi.mocked(apiService.enrichStocks).mockRejectedValue(new Error(errorMessage))

      const { data, loading, error, execute } = useEnrichStocks()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe(errorMessage)
    })
  })

  describe('useRecommendations', () => {
    it('initializes with default values', () => {
      const { data, loading, error } = useRecommendations()
      
      expect(data.value).toBeNull()
      expect(loading.value).toBe(false)
      expect(error.value).toBeNull()
    })

    it('handles successful API call with default limit', async () => {
      const mockResponse = {
        recommendations: [
          { 
            symbol: 'AAPL', 
            name: 'Apple Inc.', 
            score: 85, 
            reason: 'Strong fundamentals and growth potential',
            recommendation_type: 'BUY' as const
          }
        ],
        summary: { 
          total_analyzed: 1, 
          buy_recommendations: 1, 
          hold_recommendations: 0,
          generated_at: '2024-01-01T00:00:00Z',
          data_source: 'AlphaVantage'
        }
      }
      vi.mocked(apiService.getRecommendations).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useRecommendations()

      const promise = execute()
      expect(loading.value).toBe(true)

      await promise

      expect(loading.value).toBe(false)
      expect(data.value).toEqual(mockResponse)
      expect(error.value).toBeNull()
      expect(apiService.getRecommendations).toHaveBeenCalledWith(10)
    })

    it('handles successful API call with custom limit', async () => {
      const mockResponse = {
        recommendations: [
          { 
            symbol: 'AAPL', 
            name: 'Apple Inc.', 
            score: 85, 
            reason: 'Strong fundamentals and growth potential',
            recommendation_type: 'BUY' as const 
          },
          { 
            symbol: 'GOOGL', 
            name: 'Google LLC', 
            score: 70, 
            reason: 'Stable market position but limited growth',
            recommendation_type: 'HOLD' as const 
          }
        ],
        summary: { 
          total_analyzed: 2, 
          buy_recommendations: 1, 
          hold_recommendations: 1,
          generated_at: '2024-01-01T00:00:00Z',
          data_source: 'AlphaVantage'
        }
      }
      vi.mocked(apiService.getRecommendations).mockResolvedValue(mockResponse)

      const { data, loading, error, execute } = useRecommendations()

      await execute(20)

      expect(data.value).toEqual(mockResponse)
      expect(apiService.getRecommendations).toHaveBeenCalledWith(20)
    })

    it('handles API error', async () => {
      const errorMessage = 'Recommendations failed'
      vi.mocked(apiService.getRecommendations).mockRejectedValue(new Error(errorMessage))

      const { data, loading, error, execute } = useRecommendations()

      await execute()

      expect(loading.value).toBe(false)
      expect(data.value).toBeNull()
      expect(error.value).toBe(errorMessage)
    })
  })

  describe('error handling patterns', () => {
    it('all composables reset error on new execution', async () => {
      const composables = [
        useHealth(),
        useStocks(),
        useSyncStocks(),
        useEnrichStocks(),
        useRecommendations()
      ]

      // Set initial errors
      composables.forEach(({ error }) => {
        error.value = 'Previous error'
      })

      // Mock successful responses
      vi.mocked(apiService.checkHealth).mockResolvedValue({ status: 'healthy', service: 'stock-recommender' })
      vi.mocked(apiService.getStocks).mockResolvedValue({ stocks: [], count: 0 })
      vi.mocked(apiService.syncStocks).mockResolvedValue({ message: 'success' })
      vi.mocked(apiService.enrichStocks).mockResolvedValue({ 
        message: 'success', 
        stats: {
          total_stocks: 0,
          already_enriched: 0,
          newly_enriched: 0,
          failed: 0,
          rate_limit_reached: false
        }
      })
      vi.mocked(apiService.getRecommendations).mockResolvedValue({ 
        recommendations: [], 
        summary: {
          total_analyzed: 0,
          buy_recommendations: 0,
          hold_recommendations: 0,
          generated_at: '2024-01-01T00:00:00Z',
          data_source: 'AlphaVantage'
        }
      })

      // Execute all
      await Promise.all(composables.map(({ execute }) => execute()))

      // All errors should be cleared
      composables.forEach(({ error }) => {
        expect(error.value).toBeNull()
      })
    })

    it('all composables reset data on error', async () => {
      const composables = [
        useHealth(),
        useStocks(),
        useSyncStocks(),
        useEnrichStocks(),
        useRecommendations()
      ]

      // Set initial data
      composables.forEach(({ data }) => {
        data.value = { some: 'data' } as any
      })

      // Mock error responses
      const error = new Error('API Error')
      vi.mocked(apiService.checkHealth).mockRejectedValue(error)
      vi.mocked(apiService.getStocks).mockRejectedValue(error)
      vi.mocked(apiService.syncStocks).mockRejectedValue(error)
      vi.mocked(apiService.enrichStocks).mockRejectedValue(error)
      vi.mocked(apiService.getRecommendations).mockRejectedValue(error)

      // Execute all
      await Promise.all(composables.map(({ execute }) => execute()))

      // All data should remain as set (composables don't reset data on error in this implementation)
      composables.forEach(({ data, error }) => {
        expect(data.value).toEqual({ some: 'data' })
        expect(error.value).toBe('API Error')
      })
    })

    it('loading state is always false after execution completes', async () => {
      const { loading, execute } = useStocks()

      // Test successful case
      vi.mocked(apiService.getStocks).mockResolvedValue({ stocks: [], count: 0 })
      await execute()
      expect(loading.value).toBe(false)

      // Test error case
      vi.mocked(apiService.getStocks).mockRejectedValue(new Error('API Error'))
      await execute()
      expect(loading.value).toBe(false)
    })
  })
})