<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold text-gray-900">Smart Recommendations</h2>
      <div class="flex space-x-3">
        <select 
          v-model="selectedLimit"
          @change="handleLimitChange"
          class="border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="10">Top 10</option>
          <option value="20">Top 20</option>
          <option value="50">Top 50</option>
        </select>
        <button
          @click="handleRefresh"
          :disabled="loading.recommendations"
          class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          {{ loading.recommendations ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-md border border-gray-200 p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Filters & Sorting</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Recommendation Type</label>
          <select 
            v-model="filters.recommendation_type" 
            multiple
            class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="BUY">BUY</option>
            <option value="HOLD">HOLD</option>
            <option value="WATCH">WATCH</option>
            <option value="SELL">SELL</option>
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Sector</label>
          <select 
            v-model="filters.sector" 
            multiple
            class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option v-for="sector in availableSectors" :key="sector" :value="sector">
              {{ sector }}
            </option>
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Sort By</label>
          <select 
            v-model="sortBy"
            class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="score">Score (High to Low)</option>
            <option value="score-asc">Score (Low to High)</option>
            <option value="upside">Upside Potential</option>
            <option value="symbol">Symbol (A-Z)</option>
          </select>
        </div>
      </div>
      
      <div class="flex justify-between items-center mt-4">
        <button
          @click="clearFilters"
          class="text-sm text-gray-600 hover:text-gray-800"
        >
          Clear all filters
        </button>
        
        <span class="text-sm text-gray-500">
          Showing {{ filteredRecommendations.length }} recommendations
        </span>
      </div>
    </div>

    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400 mr-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
        </svg>
        <div class="text-sm text-red-700">{{ error }}</div>
      </div>
    </div>

    <div v-if="summary" class="bg-blue-50 border border-blue-200 rounded-md p-4">
      <h4 class="text-sm font-medium text-blue-900 mb-2">Analysis Summary</h4>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
        <div>
          <span class="text-blue-600 font-medium">Total Analyzed:</span>
          <span class="ml-1 text-blue-900">{{ summary.total_analyzed }}</span>
        </div>
        <div>
          <span class="text-blue-600 font-medium">BUY Signals:</span>
          <span class="ml-1 text-blue-900">{{ summary.buy_recommendations }}</span>
        </div>
        <div>
          <span class="text-blue-600 font-medium">HOLD Signals:</span>
          <span class="ml-1 text-blue-900">{{ summary.hold_recommendations }}</span>
        </div>
        <div>
          <span class="text-blue-600 font-medium">Data Source:</span>
          <span class="ml-1 text-blue-900">{{ summary.data_source }}</span>
        </div>
      </div>
    </div>

    <div v-if="loading.recommendations && recommendations.length === 0" class="flex items-center justify-center py-12">
      <div class="text-center">
        <svg class="animate-spin w-8 h-8 mb-2 mx-auto text-gray-400" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
        </svg>
        <p class="text-gray-500">Loading recommendations...</p>
      </div>
    </div>

    <div v-else-if="filteredRecommendations.length === 0 && !loading.recommendations" class="text-center py-12">
      <div class="text-gray-500">
        <svg class="w-16 h-16 mb-4 mx-auto text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
        </svg>
        <p class="text-lg font-medium">No recommendations found</p>
        <p class="text-sm">Try adjusting your filters or refreshing the data</p>
        <button 
          @click="clearFilters"
          class="mt-3 text-blue-600 hover:text-blue-800 text-sm underline"
        >
          Clear filters
        </button>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
      <RecommendationCard
        v-for="recommendation in filteredRecommendations"
        :key="recommendation.symbol"
        :recommendation="recommendation"
        @view-details="handleViewDetails(recommendation)"
        @add-to-watchlist="handleAddToWatchlist(recommendation)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { storeToRefs } from 'pinia'
import { useStocksStore } from '@/stores/stocks'
import RecommendationCard from './RecommendationCard.vue'
import type { Recommendation } from '@/types/api'

const store = useStocksStore()
const { recommendations, loading, error } = storeToRefs(store)

const selectedLimit = ref(10)
const sortBy = ref('score')

const filters = ref({
  recommendation_type: [] as string[],
  sector: [] as string[]
})

const summary = computed(() => {
  if (recommendations.value.length === 0) return null
  
  const buyCount = recommendations.value.filter(r => r.recommendation_type === 'BUY').length
  const holdCount = recommendations.value.filter(r => r.recommendation_type === 'HOLD').length
  
  return {
    total_analyzed: recommendations.value.length,
    buy_recommendations: buyCount,
    hold_recommendations: holdCount,
    data_source: 'Alpha Vantage + Database'
  }
})

const availableSectors = computed(() => {
  const sectors = new Set<string>()
  recommendations.value.forEach(rec => {
    if (rec.sector) sectors.add(rec.sector)
  })
  return Array.from(sectors).sort()
})

const filteredRecommendations = computed(() => {
  let filtered = recommendations.value

  if (filters.value.recommendation_type.length > 0) {
    filtered = filtered.filter(rec => 
      filters.value.recommendation_type.includes(rec.recommendation_type)
    )
  }

  if (filters.value.sector.length > 0) {
    filtered = filtered.filter(rec => 
      rec.sector && filters.value.sector.includes(rec.sector)
    )
  }

  const sorted = [...filtered].sort((a, b) => {
    switch (sortBy.value) {
      case 'score':
        return b.score - a.score
      case 'score-asc':
        return a.score - b.score
      case 'upside': {
        const aUpside = calculateUpside(a)
        const bUpside = calculateUpside(b)
        return bUpside - aUpside
      }
      case 'symbol':
        return a.symbol.localeCompare(b.symbol)
      default:
        return b.score - a.score
    }
  })

  return sorted
})

const calculateUpside = (rec: Recommendation): number => {
  if (!rec.current_price || !rec.target_price) return 0
  
  const current = parseFloat(rec.current_price)
  const target = parseFloat(rec.target_price)
  
  if (current <= 0 || target <= 0) return 0
  
  return ((target - current) / current) * 100
}

const handleRefresh = () => {
  store.fetchRecommendations(selectedLimit.value)
}

const handleLimitChange = () => {
  store.fetchRecommendations(selectedLimit.value)
}

const clearFilters = () => {
  filters.value = {
    recommendation_type: [],
    sector: []
  }
}

const handleViewDetails = (recommendation: Recommendation) => {
  console.log('View details for:', recommendation.symbol)
}

const handleAddToWatchlist = (recommendation: Recommendation) => {
  console.log('Add to watchlist:', recommendation.symbol)
}

watch(() => filters.value, () => {
}, { deep: true })

onMounted(() => {
  if (recommendations.value.length === 0) {
    handleRefresh()
  }
})
</script>