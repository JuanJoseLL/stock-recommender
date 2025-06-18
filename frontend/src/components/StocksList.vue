<template>
  <div class="max-w-6xl mx-auto p-6">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Stock Recommendations</h1>
      <p class="text-gray-600">Latest analyst recommendations and price targets</p>
    </div>

    <div v-if="error" class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
      <div class="flex items-center">
        <svg class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
        </svg>
        Error loading stocks: {{ error }}
      </div>
      <button 
        @click="retryFetch" 
        class="mt-2 text-sm underline hover:no-underline"
      >
        Try again
      </button>
    </div>

    <div class="mb-6 flex justify-between items-center">
      <div class="text-sm text-gray-500">
        Showing {{ stocksCount }} recommendations
      </div>
      <div class="flex space-x-3">
        <button 
          @click="syncStocks"
          :disabled="syncing || loading"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
        >
          <svg v-if="syncing" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
          </svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
          </svg>
          <span>{{ syncing ? 'Syncing...' : 'Sync Data' }}</span>
        </button>
        <button 
          @click="refreshStocks"
          :disabled="loading"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
        >
          <svg v-if="loading" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
          </svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          <span>{{ loading ? 'Loading...' : 'Refresh' }}</span>
        </button>
      </div>
    </div>

    <div v-if="loading && stocks.length === 0" class="text-center py-12">
      <div class="inline-flex items-center space-x-2 text-gray-500">
        <svg class="animate-spin w-6 h-6" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
        </svg>
        <span>Loading stock recommendations...</span>
      </div>
    </div>

    <div v-else-if="stocks.length === 0 && !loading" class="text-center py-12">
      <div class="text-gray-500">
        <svg class="w-12 h-12 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 48 48">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0116.07 12h15.86a2 2 0 011.664.89l.812 1.22A2 2 0 0036.07 15H37a2 2 0 012 2v18a2 2 0 01-2 2H11a2 2 0 01-2-2V17z"/>
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"/>
        </svg>
        <p class="text-lg font-medium">No stock recommendations found</p>
        <p class="text-sm">Try refreshing to load the latest data</p>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <StockCard v-for="stock in stocks" :key="`${stock.ticker}-${stock.time}-${stock.id}`" :stock="stock" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useStocksStore } from '@/stores/stocks'
import StockCard from './StockCard.vue'

const stocksStore = useStocksStore()
const { stocks, loading, error, syncing, stocksCount } = storeToRefs(stocksStore)

const refreshStocks = () => {
  stocksStore.resetStocks()
  stocksStore.fetchStocks()
}

const syncStocks = () => {
  stocksStore.syncStocks()
}

const retryFetch = () => {
  stocksStore.fetchStocks()
}

onMounted(() => {
  stocksStore.fetchStocks()
})
</script>