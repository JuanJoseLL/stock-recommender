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
          class="inline-flex items-center px-4 py-3 bg-white border border-slate-200 rounded-xl shadow-sm text-sm font-medium text-slate-700 hover:bg-slate-50 hover:border-slate-300 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 active:scale-95"
        >
          <svg class="w-4 h-4 mr-2 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          {{ loading.recommendations ? 'Loading...' : 'Refresh' }}
        </button>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-6 shadow-sm">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4 items-end">
        <!-- Search Input -->
        <div class="lg:col-span-2">
          <label for="search" class="block text-sm font-medium text-slate-700 mb-2">Search</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-4 w-4 text-slate-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              id="search"
              v-model="filters.search"
              type="text"
              placeholder="Search symbol or company..."
              class="w-full pl-10 pr-4 py-3 bg-white border border-slate-200 rounded-xl text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200"
            />
          </div>
        </div>

        <!-- Sector Filter -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Sector</label>
          <div class="relative">
            <select 
              v-model="filters.sector" 
              class="w-full appearance-none bg-white border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 cursor-pointer"
            >
              <option value="" class="text-slate-500">All Sectors</option>
              <option v-for="sector in availableSectors" :key="sector" :value="sector" class="text-slate-700">{{ sector }}</option>
            </select>
            <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
              <svg class="h-4 w-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Recommendation Type Filter -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Type</label>
          <div class="relative">
            <select 
              v-model="filters.recommendation_type" 
              class="w-full appearance-none bg-white border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 cursor-pointer"
            >
              <option value="" class="text-slate-500">All Types</option>
              <option value="BUY" class="text-slate-700">BUY</option>
              <option value="HOLD" class="text-slate-700">HOLD</option>
              <option value="WATCH" class="text-slate-700">WATCH</option>
              <option value="SELL" class="text-slate-700">SELL</option>
            </select>
            <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
              <svg class="h-4 w-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Sort By Filter -->
        <div>
          <label class="block text-sm font-medium text-slate-700 mb-2">Sort By</label>
          <div class="relative">
            <select 
              v-model="sortBy" 
              class="w-full appearance-none bg-white border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 cursor-pointer"
            >
              <option value="score" class="text-slate-700">Score (High to Low)</option>
              <option value="score-asc" class="text-slate-700">Score (Low to High)</option>
              <option value="upside" class="text-slate-700">Upside Potential</option>
              <option value="symbol" class="text-slate-700">Symbol (A-Z)</option>
            </select>
            <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
              <svg class="h-4 w-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
              </svg>
            </div>
          </div>
        </div>

        <!-- Controls -->
        <div class="flex items-center justify-between">
          <button 
            @click="clearFilters" 
            class="inline-flex items-center px-3 py-2 text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-50 border border-slate-200 rounded-xl transition-all duration-200"
          >
            <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
            Clear
          </button>
          <div class="text-right">
            <div class="text-sm font-semibold text-slate-900">{{ filteredRecommendations.length }}</div>
            <div class="text-xs text-slate-500">results</div>
          </div>
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
        <button @click="clearFilters" type="button" class="inline-flex items-center px-4 py-2 text-sm font-medium text-blue-600 hover:text-blue-500 hover:bg-blue-50 rounded-xl transition-all duration-200">
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
  sector: '',
  search: ''
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

  // Filter by search (symbol or name)
  if (filters.value.search) {
    const searchLower = filters.value.search.toLowerCase()
    filtered = filtered.filter(rec => 
      rec.symbol.toLowerCase().includes(searchLower) ||
      rec.name.toLowerCase().includes(searchLower)
    )
  }

  // Filter by recommendation type
  if (filters.value.recommendation_type) {
    filtered = filtered.filter(rec => 
      rec.recommendation_type === filters.value.recommendation_type
    )
  }

  // Filter by sector
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
    sector: '',
    search: ''
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