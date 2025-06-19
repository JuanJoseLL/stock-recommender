import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Stock {
  id?: number
  ticker: string
  target_from: string
  target_to: string
  company: string
  action: string
  brokerage: string
  rating_from: string
  rating_to: string
  time: string
  created_at?: string
  updated_at?: string
}

export interface ApiResponse {
  stocks: Stock[]
  count: number
}

export const useStocksStore = defineStore('stocks', () => {
  const stocks = ref<Stock[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const syncing = ref(false)

  const stocksCount = computed(() => stocks.value.length)

  
  const API_BASE_URL = import.meta.env.VITE_API_URL || '/api'

  const fetchStocks = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await fetch(`${API_BASE_URL}/stocks`, {
        method: 'GET',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json',
        }
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data: ApiResponse = await response.json()
      stocks.value = data.stocks || []
      
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error occurred'
      console.error('Error fetching stocks:', err)
    } finally {
      loading.value = false
    }
  }

  const syncStocks = async () => {
    syncing.value = true
    error.value = null
    
    try {
      const response = await fetch(`${API_BASE_URL}/stocks/sync`, {
        method: 'POST',
        mode: 'cors',
        headers: {
          'Content-Type': 'application/json',
        }
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      console.log('Sync response:', data.message)
      
      // Refresh stocks after sync
      await fetchStocks()
      
    } catch (err) {
      error.value = err instanceof Error ? err.message : 'Unknown error occurred'
      console.error('Error syncing stocks:', err)
    } finally {
      syncing.value = false
    }
  }

  const resetStocks = () => {
    stocks.value = []
    error.value = null
  }

  return {
    stocks,
    loading,
    error,
    syncing,
    stocksCount,
    fetchStocks,
    syncStocks,
    resetStocks
  }
})