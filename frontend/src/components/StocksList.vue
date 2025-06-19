<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="flex items-center justify-between">
       <div>
        <h1 class="text-3xl font-bold text-gray-800">Stocks</h1>
      </div>
      <div class="flex space-x-3">
        <button
          @click="handleRefresh"
          :disabled="loading.stocks"
          class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
        >
          <svg class="w-5 h-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path></svg>
          {{ loading.stocks ? 'Loading...' : 'Refresh' }}
        </button>
        <button
          @click="handleSync"
          :disabled="loading.sync"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 transition-colors"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"></path></svg>
          {{ loading.sync ? 'Syncing...' : 'Sync Data' }}
        </button>
      </div>
    </div>

    <StockFilters 
      :filters="filters" 
      :stocks="stocks"
      @update:filters="updateFilters" 
    />

    <div v-if="error" class="bg-red-50 border-l-4 border-red-400 p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9 13a1 1 0 112 0v-2a1 1 0 11-2 0v2zm0-5a1 1 0 112 0 1 1 0 01-2 0z" clip-rule="evenodd" /></svg>
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-700">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Stocks Grid -->
    <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
      <StockCard v-for="stock in paginatedStocks" :key="stock.id || stock.ticker" :stock="stock" />
    </div>

    <!-- Empty State -->
    <div v-if="filteredStocks.length === 0 && !loading.stocks" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
        <path vector-effect="non-scaling-stroke" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 13h6m-3-3v6m-9 1V7a2 2 0 012-2h6l2 2h6a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2z" />
      </svg>
      <h3 class="mt-2 text-sm font-semibold text-gray-900">No stocks found</h3>
      <p class="mt-1 text-sm text-gray-500">No stocks found matching your filters.</p>
      <div class="mt-6">
        <button @click="clearFilters" type="button" class="text-sm font-medium text-blue-600 hover:text-blue-500">
          Clear filters
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading.stocks" class="text-center py-12">
      <div class="flex items-center justify-center">
        <svg class="animate-spin w-8 h-8 mr-3 text-blue-600" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"/>
        </svg>
        <p class="text-lg text-gray-600">Loading stocks...</p>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="filteredStocks.length > pagination.limit" class="flex justify-between items-center pt-6 border-t border-gray-200">
      <div class="text-sm text-gray-700">
        Showing
        <span class="font-medium">{{ (pagination.page - 1) * pagination.limit + 1 }}</span>
        to
        <span class="font-medium">{{ Math.min(pagination.page * pagination.limit, filteredStocks.length) }}</span>
        of
        <span class="font-medium">{{ filteredStocks.length }}</span>
        results
      </div>
      <div class="flex items-center space-x-2">
        <button
          @click="previousPage"
          :disabled="pagination.page === 1"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
        >
          Previous
        </button>
        <span class="text-sm text-gray-700">
          Page {{ pagination.page }} of {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="pagination.page >= totalPages"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 transition-colors"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { storeToRefs } from 'pinia'
import { useStocksStore } from '@/stores/stocks'
import StockFilters from './StockFilters.vue'
import StockCard from './StockCard.vue'
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

const updateFilters = (newFilters: AppFilters) => {
  store.updateFilters(newFilters);
};

const clearFilters = () => {
  store.clearFilters();
};

const nextPage = () => {
  if (pagination.value.page < totalPages.value) {
    store.updatePagination({ page: pagination.value.page + 1 });
  }
};

const previousPage = () => {
  if (pagination.value.page > 1) {
    store.updatePagination({ page: pagination.value.page - 1 });
  }
};

const handleRefresh = () => {
  store.fetchStocks(true);
};

const handleSync = async () => {
  await store.syncStocks();
};

onMounted(() => {
  store.fetchStocks();
})
</script>