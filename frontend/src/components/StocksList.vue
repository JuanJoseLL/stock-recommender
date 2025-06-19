<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h2 class="text-2xl font-bold text-gray-900">Stocks List</h2>
      <div class="flex space-x-3">
        <button
          @click="handleRefresh"
          :disabled="loading.stocks"
          class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
          </svg>
          {{ loading.stocks ? 'Loading...' : 'Refresh' }}
        </button>
        <button
          @click="handleSync"
          :disabled="loading.sync"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 disabled:opacity-50"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
          </svg>
          {{ loading.sync ? 'Syncing...' : 'Sync Data' }}
        </button>
      </div>
    </div>

    <StockFilters 
      :filters="filters" 
      :stocks="stocks"
      @update:filters="updateFilters" 
    />

    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <svg class="w-5 h-5 text-red-400 mr-2" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
        </svg>
        <div class="text-sm text-red-700">{{ error }}</div>
      </div>
    </div>

    <div class="bg-white rounded-lg shadow-md border border-gray-200 overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <div class="flex justify-between items-center">
          <div class="flex items-center space-x-4">
            <span class="text-sm text-gray-700">
              Showing {{ paginatedStocks.length }} of {{ filteredStocks.length }} stocks
            </span>
            <div class="flex items-center space-x-2">
              <label class="text-sm text-gray-700">Per page:</label>
              <select 
                :value="pagination.limit"
                @change="updatePageSize"
                class="border border-gray-300 rounded-md px-2 py-1 text-sm"
              >
                <option value="10">10</option>
                <option value="20">20</option>
                <option value="50">50</option>
              </select>
            </div>
          </div>
          
          <div class="flex items-center space-x-2">
            <button
              @click="previousPage"
              :disabled="pagination.page === 1"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
            >
              Previous
            </button>
            <span class="text-sm text-gray-700">
              Page {{ pagination.page }} of {{ totalPages }}
            </span>
            <button
              @click="nextPage"
              :disabled="pagination.page >= totalPages"
              class="px-3 py-1 border border-gray-300 rounded-md text-sm disabled:opacity-50"
            >
              Next
            </button>
          </div>
        </div>
      </div>

      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Ticker
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Company
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Action
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Brokerage
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Rating
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Target Price
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Market Data
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="stock in paginatedStocks" :key="stock.id || stock.ticker" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div>
                    <div class="text-sm font-medium text-gray-900">{{ stock.ticker }}</div>
                    <div v-if="stock.sector" class="text-xs text-gray-500">{{ stock.sector }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4">
                <div class="text-sm text-gray-900">{{ stock.company }}</div>
                <div v-if="stock.industry" class="text-xs text-gray-500">{{ stock.industry }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  :class="getActionBadgeClass(stock.action)"
                >
                  {{ stock.action }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ stock.brokerage }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  {{ stock.rating_from }}
                  <span v-if="stock.rating_from !== stock.rating_to"> → {{ stock.rating_to }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  {{ stock.target_from }}
                  <span v-if="stock.target_from !== stock.target_to"> - {{ stock.target_to }}</span>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div v-if="stock.market_cap || stock.pe_ratio" class="text-sm text-gray-900">
                  <div v-if="stock.market_cap" class="text-xs">
                    Cap: {{ formatMarketCap(stock.market_cap) }}
                  </div>
                  <div v-if="stock.pe_ratio" class="text-xs">
                    P/E: {{ stock.pe_ratio }}
                  </div>
                  <div v-if="stock.week_high_52 && stock.week_low_52" class="text-xs text-gray-500">
                    52W: {{ stock.week_low_52 }} - {{ stock.week_high_52 }}
                  </div>
                </div>
                <div v-else class="text-xs text-gray-400">No data</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                  :class="stock.enriched_at ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                >
                  {{ stock.enriched_at ? 'Enriched' : 'Basic' }}
                </span>
                <div v-if="stock.enriched_at" class="text-xs text-gray-500 mt-1">
                  {{ formatDate(stock.enriched_at) }}
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="filteredStocks.length === 0 && !loading.stocks" class="px-6 py-8 text-center text-gray-500">
        <p>No stocks found matching your filters.</p>
        <button 
          @click="clearFilters"
          class="mt-2 text-blue-600 hover:text-blue-800 text-sm"
        >
          Clear filters
        </button>
      </div>

      <div v-if="loading.stocks" class="px-6 py-8 text-center text-gray-500">
        <div class="flex items-center justify-center">
          <svg class="animate-spin w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
          </svg>
          Loading stocks...
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useStocksStore } from '@/stores/stocks'
import StockFilters from './StockFilters.vue'
import type { AppFilters } from '@/types/api'

const store = useStocksStore()
const { 
  stocks, 
  loading, 
  error, 
  filters, 
  pagination, 
  filteredStocks, 
  paginatedStocks 
} = storeToRefs(store)

const totalPages = computed(() => 
  Math.ceil(filteredStocks.value.length / pagination.value.limit)
)

const getActionBadgeClass = (action: string) => {
  const actionLower = action.toLowerCase()
  if (actionLower.includes('buy') || actionLower.includes('initiated')) {
    return 'bg-green-100 text-green-800'
  } else if (actionLower.includes('sell')) {
    return 'bg-red-100 text-red-800'
  } else if (actionLower.includes('hold')) {
    return 'bg-yellow-100 text-yellow-800'
  } else {
    return 'bg-blue-100 text-blue-800'
  }
}

const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleDateString()
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

const handleRefresh = () => {
  store.fetchStocks()
}

const handleSync = () => {
  store.syncStocks()
}

const updateFilters = (newFilters: AppFilters) => {
  store.updateFilters(newFilters)
}

const updatePageSize = (event: Event) => {
  const target = event.target as HTMLSelectElement
  store.updatePagination({ 
    limit: parseInt(target.value), 
    page: 1 
  })
}

const nextPage = () => {
  if (pagination.value.page < totalPages.value) {
    store.updatePagination({ page: pagination.value.page + 1 })
  }
}

const previousPage = () => {
  if (pagination.value.page > 1) {
    store.updatePagination({ page: pagination.value.page - 1 })
  }
}

const clearFilters = () => {
  store.updateFilters({
    sector: [],
    recommendation_type: [],
    action: [],
    brokerage: [],
    enriched_only: false
  })
}

onMounted(() => {
  if (stocks.value.length === 0) {
    handleRefresh()
  }
})
</script>