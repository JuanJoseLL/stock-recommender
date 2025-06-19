import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { apiService } from '@/services/api'
import type {
  Stock,
  Recommendation,
  AppFilters,
  AppPagination,
  EnrichmentStats
} from '@/types/api'

export const useStocksStore = defineStore('stocks', () => {
  const stocks = ref<Stock[]>([])
  const recommendations = ref<Recommendation[]>([])
  const enrichmentStats = ref<EnrichmentStats | null>(null)
  
  const loading = ref({
    stocks: false,
    enrichment: false,
    recommendations: false,
    health: false,
    sync: false
  })
  
  const error = ref<string | null>(null)
  
  const filters = ref<AppFilters>({
    sector: [],
    recommendation_type: [],
    action: [],
    brokerage: [],
    enriched_only: false
  })
  
  const pagination = ref<AppPagination>({
    page: 1,
    limit: 20,
    total: 0
  })

  const stocksCount = computed(() => stocks.value.length)
  const enrichedStocksCount = computed(() => 
    stocks.value.filter(stock => stock.enriched_at).length
  )
  
  const sectorBreakdown = computed(() => {
    const sectors = stocks.value.reduce((acc, stock) => {
      if (stock.sector) {
        acc[stock.sector] = (acc[stock.sector] || 0) + 1
      }
      return acc
    }, {} as Record<string, number>)
    
    return Object.entries(sectors).map(([sector, count]) => ({ sector, count }))
  })
  
  const filteredStocks = computed(() => {
    let filtered = stocks.value
    
    if (filters.value.sector.length > 0) {
      filtered = filtered.filter(stock => 
        stock.sector && filters.value.sector.includes(stock.sector)
      )
    }
    
    if (filters.value.action.length > 0) {
      filtered = filtered.filter(stock => 
        filters.value.action.includes(stock.action)
      )
    }
    
    if (filters.value.brokerage.length > 0) {
      filtered = filtered.filter(stock => 
        filters.value.brokerage.includes(stock.brokerage)
      )
    }
    
    if (filters.value.enriched_only) {
      filtered = filtered.filter(stock => stock.enriched_at)
    }
    
    return filtered
  })
  
  const paginatedStocks = computed(() => {
    const start = (pagination.value.page - 1) * pagination.value.limit
    const end = start + pagination.value.limit
    return filteredStocks.value.slice(start, end)
  })

  const fetchStocks = async () => {
    loading.value.stocks = true
    error.value = null
    
    try {
      const response = await apiService.getStocks()
      stocks.value = response.stocks
      pagination.value.total = response.count
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch stocks'
    } finally {
      loading.value.stocks = false
    }
  }

  const syncStocks = async () => {
    loading.value.sync = true
    error.value = null
    
    try {
      await apiService.syncStocks()
      await fetchStocks()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to sync stocks'
    } finally {
      loading.value.sync = false
    }
  }

  const enrichStocks = async (limit = 5) => {
    loading.value.enrichment = true
    error.value = null
    
    try {
      const response = await apiService.enrichStocks(limit)
      enrichmentStats.value = response.stats
      await fetchStocks()
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to enrich stocks'
    } finally {
      loading.value.enrichment = false
    }
  }

  const fetchRecommendations = async (limit = 10) => {
    loading.value.recommendations = true
    error.value = null
    
    try {
      const response = await apiService.getRecommendations(limit)
      recommendations.value = response.recommendations
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Failed to fetch recommendations'
    } finally {
      loading.value.recommendations = false
    }
  }

  const updateFilters = (newFilters: Partial<AppFilters>) => {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.page = 1
  }

  const updatePagination = (newPagination: Partial<AppPagination>) => {
    pagination.value = { ...pagination.value, ...newPagination }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    stocks,
    recommendations,
    enrichmentStats,
    loading,
    error,
    filters,
    pagination,
    stocksCount,
    enrichedStocksCount,
    sectorBreakdown,
    filteredStocks,
    paginatedStocks,
    fetchStocks,
    syncStocks,
    enrichStocks,
    fetchRecommendations,
    updateFilters,
    updatePagination,
    clearError
  }
})