import type {
  HealthResponse,
  StocksResponse,
  SyncResponse,
  EnrichmentResponse,
  RecommendationsResponse
} from '@/types/api'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

class ApiService {
  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${API_BASE_URL}${endpoint}`
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      return await response.json()
    } catch (error) {
      console.error(`API request failed for ${endpoint}:`, error)
      throw error
    }
  }

  async checkHealth(): Promise<HealthResponse> {
    return this.request<HealthResponse>('/health')
  }

  async getStocks(): Promise<StocksResponse> {
    return this.request<StocksResponse>('/api/stocks')
  }

  async syncStocks(): Promise<SyncResponse> {
    return this.request<SyncResponse>('/api/stocks/sync', {
      method: 'POST',
    })
  }

  async enrichStocks(limit = 5): Promise<EnrichmentResponse> {
    return this.request<EnrichmentResponse>(`/api/stocks/enrich?limit=${limit}`, {
      method: 'POST',
    })
  }

  async getRecommendations(limit = 10): Promise<RecommendationsResponse> {
    return this.request<RecommendationsResponse>(`/api/recommendations?limit=${limit}`)
  }
}

export const apiService = new ApiService()