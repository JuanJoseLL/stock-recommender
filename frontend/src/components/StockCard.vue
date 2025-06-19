<template>
  <div class="group relative bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-6 transition-all duration-300 hover:shadow-lg hover:shadow-slate-900/5 hover:border-slate-300/50">
    <!-- Card Header -->
    <div class="flex justify-between items-start mb-6">
      <div class="flex-1">
        <div class="flex items-center space-x-3 mb-2">
          <h3 class="text-xl font-bold text-slate-900">{{ stock.ticker }}</h3>
          <span 
            class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium"
            :class="getActionBadgeClass(stock.action)"
          >
            {{ stock.action }}
          </span>
        </div>
        <p class="text-sm text-slate-600 line-clamp-1">{{ stock.company }}</p>
      </div>
    </div>

    <!-- Rating Section -->
    <div class="bg-slate-50/70 rounded-xl p-4 mb-6">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm font-medium text-slate-700">Rating Change</span>
        <div class="flex items-center space-x-2">
          <span :class="getRatingClass(stock.rating_from)" class="px-2.5 py-1 text-xs font-medium rounded-lg">
            {{ stock.rating_from }}
          </span>
          <svg class="w-4 h-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3"></path>
          </svg>
          <span :class="getRatingClass(stock.rating_to)" class="px-2.5 py-1 text-xs font-medium rounded-lg">
            {{ stock.rating_to }}
          </span>
        </div>
      </div>
      <div class="flex items-center justify-between">
        <span class="text-sm font-medium text-slate-700">Price Target</span>
        <div class="text-sm font-semibold">
          <span class="text-slate-600">${{ stock.target_from }}</span>
          <span class="text-slate-400 mx-2">→</span>
          <span class="text-emerald-600">${{ stock.target_to }}</span>
        </div>
      </div>
    </div>

    <!-- Market Data Grid -->
    <div class="grid grid-cols-2 gap-4 mb-6">
      <div class="text-center p-3 bg-gradient-to-b from-slate-50/50 to-white/50 rounded-xl border border-slate-100">
        <p class="text-xs text-slate-500 mb-1">Market Cap</p>
        <p class="text-sm font-semibold text-slate-900">{{ formatMarketCap(stock.market_cap) }}</p>
      </div>
      <div class="text-center p-3 bg-gradient-to-b from-slate-50/50 to-white/50 rounded-xl border border-slate-100">
        <p class="text-xs text-slate-500 mb-1">P/E Ratio</p>
        <p class="text-sm font-semibold text-slate-900">{{ stock.pe_ratio || 'N/A' }}</p>
      </div>
    </div>

    <!-- 52-Week Range -->
    <div class="mb-6">
      <div class="flex justify-between items-center mb-2">
        <span class="text-xs text-slate-500">52-Week Range</span>
        <span class="text-xs font-medium text-slate-700">{{ stock.week_low_52 }} - {{ stock.week_high_52 }}</span>
      </div>
      <div class="w-full bg-slate-200 rounded-full h-1.5">
        <div class="bg-gradient-to-r from-blue-500 to-indigo-500 h-1.5 rounded-full" style="width: 65%"></div>
      </div>
    </div>

    <!-- Card Footer -->
    <div class="flex items-center justify-between pt-4 border-t border-slate-100">
      <div class="flex items-center space-x-3">
        <div class="w-8 h-8 bg-gradient-to-br from-blue-100 to-indigo-100 rounded-full flex items-center justify-center">
          <div class="w-3 h-3 bg-blue-600 rounded-full"></div>
        </div>
        <div>
          <p class="text-sm font-medium text-slate-900">{{ stock.brokerage }}</p>
          <p class="text-xs text-slate-500">{{ stock.sector }}</p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <div 
          class="w-2 h-2 rounded-full"
          :class="stock.enriched_at ? 'bg-emerald-400' : 'bg-slate-300'"
        ></div>
        <span 
          class="text-xs font-medium"
          :class="stock.enriched_at ? 'text-emerald-700' : 'text-slate-500'"
        >
          {{ stock.enriched_at ? 'Enriched' : 'Basic' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Stock } from '@/types/api'

interface Props {
  stock: Stock
}

const props = defineProps<Props>()

const getActionBadgeClass = (action: string) => {
  const actionLower = action?.toLowerCase() || ''
  if (actionLower.includes('buy') || actionLower.includes('initiated') || actionLower.includes('upgrade')) {
    return 'bg-green-100 text-green-800'
  } else if (actionLower.includes('sell') || actionLower.includes('downgrade')) {
    return 'bg-red-100 text-red-800'
  } else if (actionLower.includes('hold')) {
    return 'bg-yellow-100 text-yellow-800'
  } else {
    return 'bg-blue-100 text-blue-800'
  }
}

const getRatingClass = (rating: string) => {
  const r = rating?.toLowerCase() || ''
  if (r.includes('buy') || r.includes('outperform')) {
    return 'bg-green-100 text-green-800'
  } else if (r.includes('sell') || r.includes('underperform')) {
    return 'bg-red-100 text-red-800'
  } else {
    return 'bg-gray-200 text-gray-800'
  }
}

const formatMarketCap = (marketCap?: number) => {
  if (!marketCap) return 'N/A'
  if (marketCap >= 1e12) {
    return `$${(marketCap / 1e12).toFixed(2)}T`
  } else if (marketCap >= 1e9) {
    return `$${(marketCap / 1e9).toFixed(2)}B`
  } else if (marketCap >= 1e6) {
    return `$${(marketCap / 1e6).toFixed(2)}M`
  }
  return `$${marketCap}`
}
</script>