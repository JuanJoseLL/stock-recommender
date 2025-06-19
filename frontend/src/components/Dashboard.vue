<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-slate-900 tracking-tight">Dashboard</h1>
      </div>
      <div class="flex space-x-3">
        <button
          @click="refreshData"

          :disabled="loading.stocks || loading.recommendations"
          class="inline-flex items-center px-4 py-3 bg-white border border-slate-200 rounded-xl shadow-sm text-sm font-medium text-slate-700 hover:bg-slate-50 hover:border-slate-300 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 active:scale-95"
        >
          <svg class="w-4 h-4 mr-2 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          Refresh
        </button>
        <button
          @click="handleEnrichData"
          :disabled="loading.enrichment"
          class="inline-flex items-center px-4 py-3 bg-slate-900 hover:bg-slate-800 text-white rounded-xl shadow-sm text-sm font-medium disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 active:scale-95"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
          </svg>
          {{ loading.enrichment ? 'Enriching...' : 'Enrich Data' }}
        </button>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-50/80 backdrop-blur-sm border border-red-200 rounded-xl p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9 13a1 1 0 112 0v-2a1 1 0 11-2 0v2zm0-5a1 1 0 112 0 1 1 0 01-2 0z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Dashboard Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 mb-16">
      <DashboardCard
        title="Total Stocks"
        :value="stocksCount"    
        icon="chart"
        
        subtitle="Tracked securities"
      />
      <DashboardCard
        title="Enriched Data"
        :value="enrichedStocksCount"
        icon="check"
        :subtitle="`${enrichmentPercentage}% complete`"
      />
      <DashboardCard
        title="Top Recommendations"
        :value="buyRecommendationsCount"
        icon="target"
        subtitle="BUY signals"
      />
      <DashboardCard
        title="Sectors Tracked"
        :value="sectorBreakdown.length"
        icon="building"
        subtitle="Market sectors"
      />
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-5 gap-8 mt-5">
      <div class="lg:col-span-3 bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-8 shadow-sm mt-8">
        <h3 class="text-xl font-semibold text-slate-900 mb-8">Sector Breakdown</h3>
        <div v-if="sectorBreakdown.length > 0" class="space-y-6">
          <div 
            v-for="sector in sectorBreakdown.slice(0, 8)" 
            :key="sector.sector"
          >
            <div class="flex items-center justify-between mb-3">
              <span class="text-sm font-medium text-slate-700">{{ sector.sector }}</span>
              <span class="text-sm font-semibold text-slate-900">{{ sector.count }}</span>
            </div>
            <div class="w-full bg-slate-100 rounded-full h-2.5">
              <div 
                class="bg-gradient-to-r from-blue-500 to-indigo-500 h-2.5 rounded-full" 
                :style="{ width: `${(sector.count / stocksCount) * 100}%` }"
              ></div>
            </div>
          </div>
        </div>
        <div v-else class="text-center text-slate-500 py-16">
          <svg class="w-16 h-16 mx-auto mb-6 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"></path>
          </svg>
          <p class="text-sm">No sector data available</p>
        </div>
      </div>

      <div class="lg:col-span-2 bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-8 shadow-sm">
        <h3 class="text-xl font-semibold text-slate-900 mb-8">Top Recommendations</h3>
        <div v-if="recommendations.length > 0" class="space-y-4">
          <div 
            v-for="rec in recommendations.slice(0, 5)" 
            :key="rec.symbol"
            class="p-5 rounded-xl border border-slate-200 bg-slate-50/30"
          >
            <div class="flex items-center justify-between mb-3">
              <span class="font-bold text-slate-900 text-lg">{{ rec.symbol }}</span>
              <span 
                class="inline-flex items-center px-3 py-1.5 rounded-full text-xs font-medium"
                :class="getRecommendationBadgeColor(rec.recommendation_type)"
              >
                {{ rec.recommendation_type }}
              </span>
            </div>
            <p class="text-sm text-slate-600 line-clamp-1">{{ rec.name }}</p>
          </div>
        </div>
        <div v-else class="text-center text-slate-500 py-16">
          <svg class="w-16 h-16 mx-auto mb-6 text-slate-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"></path>
          </svg>
          <p class="text-sm">No recommendations available</p>
        </div>
      </div>
    </div>

    <!-- Enrichment Stats -->
    <div v-if="enrichmentStats" class="bg-gradient-to-r from-blue-50/80 to-indigo-50/80 backdrop-blur-sm border border-blue-200/50 rounded-2xl p-6 shadow-sm">
      <div class="flex items-center space-x-3 mb-4">
        <div class="w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center">
          <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
          </svg>
        </div>
        <h4 class="text-lg font-semibold text-slate-900">Latest Enrichment Results</h4>
      </div>
      <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
        <div class="text-center p-4 bg-white/50 rounded-xl border border-blue-100">
          <p class="text-sm font-medium text-blue-600 mb-1">Total</p>
          <p class="text-2xl font-bold text-slate-900">{{ enrichmentStats.total_stocks }}</p>
        </div>
        <div class="text-center p-4 bg-white/50 rounded-xl border border-emerald-100">
          <p class="text-sm font-medium text-emerald-600 mb-1">Enriched</p>
          <p class="text-2xl font-bold text-slate-900">{{ enrichmentStats.newly_enriched }}</p>
        </div>
        <div class="text-center p-4 bg-white/50 rounded-xl border border-slate-100">
          <p class="text-sm font-medium text-slate-600 mb-1">Already Done</p>
          <p class="text-2xl font-bold text-slate-900">{{ enrichmentStats.already_enriched }}</p>
        </div>
        <div class="text-center p-4 bg-white/50 rounded-xl border border-red-100">
          <p class="text-sm font-medium text-red-600 mb-1">Failed</p>
          <p class="text-2xl font-bold text-slate-900">{{ enrichmentStats.failed }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useStocksStore } from '@/stores/stocks'
import DashboardCard from './DashboardCard.vue'

const store = useStocksStore()
const { 
  stocksCount, 
  enrichedStocksCount, 
  sectorBreakdown, 
  recommendations,
  enrichmentStats,
  loading,
  error 
} = storeToRefs(store)

const enrichmentPercentage = computed(() => {
  if (stocksCount.value === 0) return 0
  return Math.round((enrichedStocksCount.value / stocksCount.value) * 100)
})

const buyRecommendationsCount = computed(() => {
  return recommendations.value.filter(rec => rec.recommendation_type === 'BUY').length
})

const getRecommendationBgColor = (type: string) => {
  switch (type) {
    case 'BUY': return 'bg-green-50'
    case 'HOLD': return 'bg-yellow-50'
    case 'WATCH': return 'bg-blue-50'
    case 'SELL': return 'bg-red-50'
    default: return 'bg-gray-50'
  }
}

const getRecommendationBadgeColor = (type: string) => {
  switch (type) {
    case 'BUY': return 'bg-green-100 text-green-800'
    case 'HOLD': return 'bg-yellow-100 text-yellow-800'
    case 'WATCH': return 'bg-blue-100 text-blue-800'
    case 'SELL': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const refreshData = async () => {
  await Promise.all([
    store.fetchStocks(),
    store.fetchRecommendations(10)
  ])
}

const handleEnrichData = async () => {
  await store.enrichStocks(10)
}

onMounted(() => {
  refreshData()
})
</script>