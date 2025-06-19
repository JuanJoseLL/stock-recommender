import { ref, type Ref } from 'vue'
import { apiService } from '@/services/api'
import type {
  HealthResponse,
  StocksResponse,
  SyncResponse,
  EnrichmentResponse,
  RecommendationsResponse
} from '@/types/api'

interface UseApiReturn<T> {
  data: Ref<T | null>
  loading: Ref<boolean>
  error: Ref<string | null>
  execute: (...args: any[]) => Promise<void>
}

export function useHealth(): UseApiReturn<HealthResponse> {
  const data = ref<HealthResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const execute = async () => {
    loading.value = true
    error.value = null
    try {
      data.value = await apiService.checkHealth()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to check health'
    } finally {
      loading.value = false
    }
  }

  return { data, loading, error, execute }
}

export function useStocks(): UseApiReturn<StocksResponse> {
  const data = ref<StocksResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const execute = async () => {
    loading.value = true
    error.value = null
    try {
      data.value = await apiService.getStocks()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stocks'
    } finally {
      loading.value = false
    }
  }

  return { data, loading, error, execute }
}

export function useSyncStocks(): UseApiReturn<SyncResponse> {
  const data = ref<SyncResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const execute = async () => {
    loading.value = true
    error.value = null
    try {
      data.value = await apiService.syncStocks()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to sync stocks'
    } finally {
      loading.value = false
    }
  }

  return { data, loading, error, execute }
}

export function useEnrichStocks(): UseApiReturn<EnrichmentResponse> {
  const data = ref<EnrichmentResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const execute = async (limit = 5) => {
    loading.value = true
    error.value = null
    try {
      data.value = await apiService.enrichStocks(limit)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to enrich stocks'
    } finally {
      loading.value = false
    }
  }

  return { data, loading, error, execute }
}

export function useRecommendations(): UseApiReturn<RecommendationsResponse> {
  const data = ref<RecommendationsResponse | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const execute = async (limit = 10) => {
    loading.value = true
    error.value = null
    try {
      data.value = await apiService.getRecommendations(limit)
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch recommendations'
    } finally {
      loading.value = false
    }
  }

  return { data, loading, error, execute }
}