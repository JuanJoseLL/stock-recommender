<template>
  <div class="group relative bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-6 flex flex-col justify-between transition-all duration-300 hover:shadow-lg hover:shadow-slate-900/5 hover:border-slate-300/50">
    <div>
      <!-- Card Header -->
      <div class="flex justify-between items-start mb-6">
        <div class="flex-1">
          <div class="flex items-center space-x-3 mb-2">
            <h3 class="text-xl font-bold text-slate-900">{{ recommendation.symbol }}</h3>
            <span 
              class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium"
              :class="getRecommendationBadgeColor(recommendation.recommendation_type)"
            >
              {{ recommendation.recommendation_type }}
            </span>
          </div>
          <p class="text-sm text-slate-600 line-clamp-1">{{ recommendation.name }}</p>
        </div>
      </div>

      <!-- Score Section -->
      <div class="bg-gradient-to-r from-slate-50/70 to-blue-50/70 rounded-xl p-4 mb-6">
        <div class="flex items-center justify-between mb-3">
          <span class="text-sm font-medium text-slate-700">Recommendation Score</span>
          <div class="text-right">
            <span class="text-2xl font-bold" :class="getScoreTextColor(recommendation.score)">
              {{ recommendation.score }}
            </span>
            <span class="text-sm text-slate-500">/100</span>
          </div>
        </div>
        <div class="w-full bg-slate-200 rounded-full h-2">
          <div 
            class="h-2 rounded-full transition-all duration-500 ease-out" 
            :class="getScoreBarColor(recommendation.score)" 
            :style="{ width: `${recommendation.score}%` }"
          ></div>
        </div>
      </div>

      <!-- Reason -->
      <div class="mb-6">
        <p class="text-sm text-slate-600 leading-relaxed line-clamp-3">{{ recommendation.reason }}</p>
      </div>
      
      <!-- Price Data -->
      <div class="grid grid-cols-2 gap-3 mb-6">
        <div class="text-center p-4 bg-gradient-to-b from-slate-50/50 to-white/50 rounded-xl border border-slate-100">
          <p class="text-xs text-slate-500 mb-2">Current Price</p>
          <p class="text-lg font-bold text-slate-900">${{ recommendation.current_price }}</p>
        </div>
        <div class="text-center p-4 bg-gradient-to-b from-slate-50/50 to-white/50 rounded-xl border border-slate-100">
          <p class="text-xs text-slate-500 mb-2">Target Price</p>
          <p class="text-lg font-bold" :class="upside && upside.value > 0 ? 'text-emerald-600' : 'text-red-500'">
            ${{ recommendation.target_price }}
          </p>
        </div>
      </div>

      <!-- Upside Potential -->
      <div v-if="upside" class="p-4 rounded-xl mb-6" :class="upside.bgColor">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium" :class="upside.textColor">Upside Potential</span>
          <span class="text-lg font-bold" :class="upside.textColor">{{ upside.formattedValue }}</span>
        </div>
      </div>
    </div>

    <!-- Card Footer Actions -->
    <div class="flex space-x-3 pt-4 border-t border-slate-100">
      <button 
        @click="$emit('view-details')"
        class="flex-1 inline-flex justify-center items-center px-4 py-3 bg-slate-900 hover:bg-slate-800 text-white text-sm font-medium rounded-xl transition-all duration-200 active:scale-95"
      >
        View Details
      </button>
      <button 
        @click="$emit('add-to-watchlist')"
        class="flex-1 inline-flex justify-center items-center px-4 py-3 bg-slate-100 hover:bg-slate-200 text-slate-700 text-sm font-medium rounded-xl transition-all duration-200 active:scale-95"
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
  
  const isPositive = upsidePercent > 0;
  
  return {
    value: upsidePercent,
    formattedValue: `${isPositive ? '+' : ''}${upsidePercent.toFixed(1)}%`,
    bgColor: isPositive ? 'bg-green-50' : 'bg-red-50',
    textColor: isPositive ? 'text-green-700' : 'text-red-700'
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

const getScoreTextColor = (score: number) => {
  if (score >= 80) return 'text-green-600'
  if (score >= 60) return 'text-yellow-600'
  if (score >= 40) return 'text-orange-600'
  return 'text-red-600'
}
</script>