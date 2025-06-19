<template>
  <div class="space-y-8">
        <!-- Header -->
    <div class="flex items-center justify-between">
       <div>
        <h1 class="text-3xl font-bold text-gray-800">Recommendations</h1>
      </div>
      <div class="flex space-x-3">
        <button
          @click="handleRefresh"
          :disabled="loading.recommendations"
          class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
        >
          <svg class="w-5 h-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          {{ loading.recommendations ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white rounded-xl border border-gray-200 p-4">
      <div class="flex flex-col md:flex-row items-center space-y-4 md:space-y-0 md:space-x-4">
        <div class="w-full md:w-auto">
          <label class="sr-only" for="rec-type">Recommendation Type</label>
          <select id="rec-type" v-model="filters.recommendation_type" class="w-full border-gray-300 rounded-lg shadow-sm focus:border-blue-500 focus:ring-blue-500">
            <option value="">All Types</option>
            <option value="BUY">BUY</option>
            <option value="HOLD">HOLD</option>
            <option value="WATCH">WATCH</option>
            <option value="SELL">SELL</option>
          </select>
        </div>
        <div class="w-full md:w-auto">
          <label class="sr-only" for="sector">Sector</label>
          <select id="sector" v-model="filters.sector" class="w-full border-gray-300 rounded-lg shadow-sm focus:border-blue-500 focus:ring-blue-500">
            <option value="">All Sectors</option>
            <option v-for="sector in availableSectors" :key="sector" :value="sector">{{ sector }}</option>
          </select>
        </div>
        <div class="w-full md:w-auto">
          <label class="sr-only" for="sort-by">Sort By</label>
          <select id="sort-by" v-model="sortBy" class="w-full border-gray-300 rounded-lg shadow-sm focus:border-blue-500 focus:ring-blue-500">
            <option value="score">Score (High to Low)</option>
            <option value="score-asc">Score (Low to High)</option>
            <option value="upside">Upside Potential</option>
            <option value="symbol">Symbol (A-Z)</option>
          </select>
        </div>
        <button @click="clearFilters" class="text-sm text-blue-600 hover:text-blue-800 transition-colors">
          Clear Filters
        </button>
        <div class="text-sm text-gray-500 md:ml-auto">
          Showing {{ filteredRecommendations.length }} recommendations
        </div>
      </div>
    </div>

    <!-- Summary -->
    <div v-if="summary" class="bg-blue-50 border-l-4 border-blue-400 p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-blue-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" /></svg>
        </div>
        <div class="ml-3 flex-1 md:flex md:justify-between">
          <p class="text-sm text-blue-700">
            Analyzed <span class="font-medium">{{ summary.total_analyzed }}</span> stocks. Found <span class="font-medium">{{ summary.buy_recommendations }}</span> BUY and <span class="font-medium">{{ summary.hold_recommendations }}</span> HOLD signals.
          </p>
          <p class="mt-3 text-sm md:mt-0 md:ml-6">
            <span class="text-blue-700">Source: <span class="font-medium">{{ summary.data_source }}</span></span>
          </p>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading.recommendations && recommendations.length === 0" class="text-center py-12">
      <div class="flex items-center justify-center">
        <svg class="animate-spin w-8 h-8 mr-3 text-blue-600" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/></svg>
        <p class="text-lg text-gray-600">Loading recommendations...</p>
      </div>
    </div>
    
    <!-- Empty State -->
    <div v-else-if="filteredRecommendations.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true"><path vector-effect="non-scaling-stroke" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" /></svg>
      <h3 class="mt-2 text-sm font-semibold text-gray-900">No recommendations found</h3>
      <p class="mt-1 text-sm text-gray-500">Try adjusting your filters or refreshing the data.</p>
      <div class="mt-6">
        <button @click="clearFilters" type="button" class="text-sm font-medium text-blue-600 hover:text-blue-500">
          Clear filters
        </button>
      </div>
    </div>

    <!-- Recommendations Grid -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
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
  recommendation_type: '',
  sector: ''
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

  if (filters.value.recommendation_type) {
    filtered = filtered.filter(rec => 
      rec.recommendation_type === filters.value.recommendation_type
    )
  }

  if (filters.value.sector) {
    filtered = filtered.filter(rec => 
      rec.sector && rec.sector === filters.value.sector
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
    recommendation_type: '',
    sector: ''
  }
}

const handleViewDetails = (rec: Recommendation) => {
  // Logic to show details, e.g., in a modal
  alert(`Viewing details for ${rec.symbol}`)
}

const handleAddToWatchlist = (rec: Recommendation) => {
  // Logic to add to a watchlist
  alert(`${rec.symbol} added to watchlist`)
}

watch(() => filters.value, () => {
}, { deep: true })

onMounted(() => {
  if (recommendations.value.length === 0) {
    handleRefresh()
  }
})
</script>