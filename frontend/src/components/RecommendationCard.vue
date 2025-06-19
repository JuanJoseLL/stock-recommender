<template>
  <div class="bg-white rounded-lg shadow-md border border-gray-200 p-6 hover:shadow-lg transition-shadow">
    <div class="flex items-start justify-between mb-4">
      <div class="flex-1">
        <div class="flex items-center space-x-3">
          <h3 class="text-lg font-semibold text-gray-900">{{ recommendation.symbol }}</h3>
          <span 
            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
            :class="getRecommendationBadgeColor(recommendation.recommendation_type)"
          >
            {{ recommendation.recommendation_type }}
          </span>
        </div>
        <p class="text-sm text-gray-600 mt-1">{{ recommendation.name }}</p>
        <p v-if="recommendation.sector" class="text-xs text-gray-500 mt-1">{{ recommendation.sector }}</p>
      </div>
      
      <div class="text-right">
        <div class="text-2xl font-bold text-gray-900">{{ recommendation.score }}</div>
        <div class="text-xs text-gray-500">/ 100</div>
      </div>
    </div>

    <div class="mb-4">
      <div class="flex items-center justify-between text-sm text-gray-600 mb-2">
        <span>Score</span>
        <span>{{ recommendation.score }}/100</span>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-2">
        <div 
          class="h-2 rounded-full transition-all duration-300"
          :class="getScoreBarColor(recommendation.score)"
          :style="{ width: `${recommendation.score}%` }"
        ></div>
      </div>
    </div>

    <div v-if="recommendation.current_price || recommendation.target_price" class="grid grid-cols-2 gap-4 mb-4">
      <div v-if="recommendation.current_price" class="text-center p-3 bg-gray-50 rounded-lg">
        <div class="text-xs text-gray-500">Current Price</div>
        <div class="text-lg font-semibold text-gray-900">${{ recommendation.current_price }}</div>
      </div>
      <div v-if="recommendation.target_price" class="text-center p-3 bg-blue-50 rounded-lg">
        <div class="text-xs text-gray-500">Target Price</div>
        <div class="text-lg font-semibold text-blue-900">${{ recommendation.target_price }}</div>
      </div>
    </div>

    <div v-if="upside" class="mb-4 p-3 rounded-lg" :class="upside.color">
      <div class="flex items-center justify-between">
        <span class="text-sm font-medium">Upside Potential</span>
        <span class="text-sm font-bold">{{ upside.value }}</span>
      </div>
    </div>

    <div v-if="recommendation.pe_ratio || recommendation.market_cap" class="grid grid-cols-2 gap-4 mb-4 text-sm">
      <div v-if="recommendation.pe_ratio" class="flex justify-between">
        <span class="text-gray-500">P/E Ratio:</span>
        <span class="font-medium">{{ recommendation.pe_ratio }}</span>
      </div>
      <div v-if="recommendation.market_cap" class="flex justify-between">
        <span class="text-gray-500">Market Cap:</span>
        <span class="font-medium">{{ formatMarketCap(Number(recommendation.market_cap)) }}</span>
      </div>
    </div>

    <div class="border-t border-gray-200 pt-4">
      <p class="text-sm text-gray-700 leading-relaxed">{{ recommendation.reason }}</p>
    </div>

    <div class="mt-4 flex space-x-2">
      <button 
        @click="$emit('view-details')"
        class="flex-1 px-4 py-2 text-sm font-medium text-blue-600 bg-blue-50 border border-blue-200 rounded-md hover:bg-blue-100 transition-colors"
      >
        View Details
      </button>
      <button 
        @click="$emit('add-to-watchlist')"
        class="px-4 py-2 text-sm font-medium text-gray-600 bg-gray-50 border border-gray-200 rounded-md hover:bg-gray-100 transition-colors"
      >
        Add to Watchlist
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Recommendation } from '@/types/api'

interface Props {
  recommendation: Recommendation
}

interface Emits {
  (e: 'view-details'): void
  (e: 'add-to-watchlist'): void
}

const props = defineProps<Props>()
defineEmits<Emits>()

const upside = computed(() => {
  if (!props.recommendation.current_price || !props.recommendation.target_price) {
    return null
  }
  
  const current = parseFloat(props.recommendation.current_price)
  const target = parseFloat(props.recommendation.target_price)
  
  if (current <= 0 || target <= 0) return null
  
  const upsidePercent = ((target - current) / current) * 100
  
  if (upsidePercent > 0) {
    return {
      value: `+${upsidePercent.toFixed(1)}%`,
      color: 'bg-green-50 text-green-700 border border-green-200'
    }
  } else {
    return {
      value: `${upsidePercent.toFixed(1)}%`,
      color: 'bg-red-50 text-red-700 border border-red-200'
    }
  }
})

const getRecommendationBadgeColor = (type: string) => {
  switch (type) {
    case 'BUY': return 'bg-green-100 text-green-800'
    case 'HOLD': return 'bg-yellow-100 text-yellow-800'
    case 'WATCH': return 'bg-blue-100 text-blue-800'
    case 'SELL': return 'bg-red-100 text-red-800'
    default: return 'bg-gray-100 text-gray-800'
  }
}

const getScoreBarColor = (score: number) => {
  if (score >= 80) return 'bg-green-500'
  if (score >= 60) return 'bg-yellow-500'
  if (score >= 40) return 'bg-orange-500'
  return 'bg-red-500'
}

const formatMarketCap = (marketCap: number) => {
  if (marketCap >= 1e9) {
    return `$${(marketCap / 1e9).toFixed(1)}B`
  } else if (marketCap >= 1e6) {
    return `$${(marketCap / 1e6).toFixed(1)}M`
  } else {
    return `$${marketCap.toLocaleString()}`
  }
}
</script>