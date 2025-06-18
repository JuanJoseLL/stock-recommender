<template>
  <div class="bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow p-6">
    <div class="flex justify-between items-start mb-4">
      <div>
        <h3 class="text-lg font-semibold text-gray-900">{{ stock.ticker }}</h3>
        <p class="text-sm text-gray-600">{{ stock.company }}</p>
      </div>
      <span :class="actionBadgeClass" class="px-2 py-1 text-xs font-medium rounded-full">
        {{ stock.action }}
      </span>
    </div>
    
    <div class="space-y-3">
      <div class="flex justify-between">
        <span class="text-sm text-gray-500">Rating</span>
        <div class="flex items-center space-x-2">
          <span :class="ratingFromClass" class="px-2 py-1 text-xs font-medium rounded">
            {{ stock.rating_from }}
          </span>
          <span class="text-gray-400">→</span>
          <span :class="ratingToClass" class="px-2 py-1 text-xs font-medium rounded">
            {{ stock.rating_to }}
          </span>
        </div>
      </div>
      
      <div class="flex justify-between">
        <span class="text-sm text-gray-500">Price Target</span>
        <div class="flex items-center space-x-2">
          <span class="text-sm font-medium text-gray-700">{{ stock.target_from }}</span>
          <span class="text-gray-400">→</span>
          <span class="text-sm font-medium text-green-600">{{ stock.target_to }}</span>
        </div>
      </div>
      
      <div class="flex justify-between items-center pt-3 border-t border-gray-100">
        <span class="text-sm text-gray-500">{{ stock.brokerage }}</span>
        <span class="text-xs text-gray-400">{{ formattedTime }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Stock } from '@/stores/stocks'

interface Props {
  stock: Stock
}

const props = defineProps<Props>()

const formattedTime = computed(() => {
  return new Date(props.stock.time).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric'
  })
})

const actionBadgeClass = computed(() => {
  const action = props.stock.action.toLowerCase()
  if (action.includes('upgrade')) {
    return 'bg-green-100 text-green-800'
  } else if (action.includes('downgrade')) {
    return 'bg-red-100 text-red-800'
  } else {
    return 'bg-blue-100 text-blue-800'
  }
})

const getRatingClass = (rating: string) => {
  const r = rating.toLowerCase()
  if (r.includes('buy') || r.includes('outperform')) {
    return 'bg-green-100 text-green-800'
  } else if (r.includes('sell') || r.includes('underperform')) {
    return 'bg-red-100 text-red-800'
  } else {
    return 'bg-gray-100 text-gray-800'
  }
}

const ratingFromClass = computed(() => getRatingClass(props.stock.rating_from))
const ratingToClass = computed(() => getRatingClass(props.stock.rating_to))
</script>