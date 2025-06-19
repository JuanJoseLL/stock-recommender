<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-2xl font-bold text-gray-900">Stock Recommender Dashboard</h1>
      <div class="flex space-x-3">
        <button
          @click="refreshData"
          :disabled="loading.stocks || loading.recommendations"
          class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          Refresh
        </button>
        <button
          @click="handleEnrichData"
          :disabled="loading.enrichment"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
          </svg>
          {{ loading.enrichment ? 'Enriching...' : 'Enrich Data' }}
        </button>
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

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <DashboardCard
        title="Total Stocks"
        :value="stocksCount"
        icon="chart"
        iconColor="text-blue-600"
        iconBgColor="bg-blue-100"
        subtitle="Tracked securities"
      />
      
      <DashboardCard
        title="Enriched Data"
        :value="enrichedStocksCount"
        icon="check"
        iconColor="text-green-600"
        iconBgColor="bg-green-100"
        :subtitle="`${enrichmentPercentage}% complete`"
      />
      
      <DashboardCard
        title="Top Recommendations"
        :value="buyRecommendationsCount"
        icon="target"
        iconColor="text-purple-600"
        iconBgColor="bg-purple-100"
        subtitle="BUY signals"
      />
      
      <DashboardCard
        title="Sectors Tracked"
        :value="sectorBreakdown.length"
        icon="building"
        iconColor="text-indigo-600"
        iconBgColor="bg-indigo-100"
        subtitle="Market sectors"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white rounded-lg shadow-md p-6 border border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Sector Breakdown</h3>
        <div v-if="sectorBreakdown.length > 0" class="space-y-3">
          <div 
            v-for="sector in sectorBreakdown.slice(0, 8)" 
            :key="sector.sector"
            class="flex items-center justify-between"
          >
            <span class="text-sm font-medium text-gray-700">{{ sector.sector }}</span>
            <div class="flex items-center">
              <div class="w-20 bg-gray-200 rounded-full h-2 mr-3">
                <div 
                  class="bg-blue-600 h-2 rounded-full" 
                  :style="{ width: `${(sector.count / stocksCount) * 100}%` }"
                ></div>
              </div>
              <span class="text-sm text-gray-500 w-8 text-right">{{ sector.count }}</span>
            </div>
          </div>
        </div>
        <div v-else class="text-center text-gray-500 py-8">
          No sector data available
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6 border border-gray-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Top Recommendations</h3>
        <div v-if="recommendations.length > 0" class="space-y-3">
          <div 
            v-for="rec in recommendations.slice(0, 5)" 
            :key="rec.symbol"
            class="flex items-center justify-between p-3 rounded-lg"
            :class="getRecommendationBgColor(rec.recommendation_type)"
          >
            <div class="flex-1">
              <div class="flex items-center justify-between">
                <span class="font-medium text-gray-900">{{ rec.symbol }}</span>
                <span 
                  class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                  :class="getRecommendationBadgeColor(rec.recommendation_type)"
                >
                  {{ rec.recommendation_type }}
                </span>
              </div>
              <p class="text-sm text-gray-600 mt-1">{{ rec.name }}</p>
              <div class="flex items-center mt-2">
                <div class="w-16 bg-gray-200 rounded-full h-1.5 mr-2">
                  <div 
                    class="bg-green-500 h-1.5 rounded-full" 
                    :style="{ width: `${rec.score}%` }"
                  ></div>
                </div>
                <span class="text-xs text-gray-500">{{ rec.score }}/100</span>
              </div>
            </div>
          </div>
        </div>
        <div v-else class="text-center text-gray-500 py-8">
          No recommendations available
        </div>
      </div>
    </div>

    <div v-if="enrichmentStats" class="bg-blue-50 border border-blue-200 rounded-md p-4">
      <h4 class="text-sm font-medium text-blue-900 mb-2">Latest Enrichment Results</h4>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
        <div>
          <span class="text-blue-600 font-medium">Total:</span>
          <span class="ml-1 text-blue-900">{{ enrichmentStats.total_stocks }}</span>
        </div>
        <div>
          <span class="text-blue-600 font-medium">Enriched:</span>
          <span class="ml-1 text-blue-900">{{ enrichmentStats.newly_enriched }}</span>
        </div>
        <div>
          <span class="text-blue-600 font-medium">Already Done:</span>
          <span class="ml-1 text-blue-900">{{ enrichmentStats.already_enriched }}</span>
        </div>
        <div>
          <span class="text-blue-600 font-medium">Failed:</span>
          <span class="ml-1 text-blue-900">{{ enrichmentStats.failed }}</span>
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