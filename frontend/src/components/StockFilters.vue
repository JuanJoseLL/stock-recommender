<template>
  <div class="bg-white rounded-lg shadow-md p-6 border border-gray-200">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Filters</h3>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Sector</label>
        <select 
          v-model="localFilters.sector" 
          multiple
          class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option v-for="sector in availableSectors" :key="sector" :value="sector">
            {{ sector }}
          </option>
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Action</label>
        <select 
          v-model="localFilters.action" 
          multiple
          class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option v-for="action in availableActions" :key="action" :value="action">
            {{ action }}
          </option>
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Brokerage</label>
        <select 
          v-model="localFilters.brokerage" 
          multiple
          class="w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option v-for="brokerage in availableBrokerages" :key="brokerage" :value="brokerage">
            {{ brokerage }}
          </option>
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
        <div class="space-y-2">
          <label class="flex items-center">
            <input 
              type="checkbox" 
              v-model="localFilters.enriched_only"
              class="mr-2 text-blue-600 focus:ring-blue-500"
            >
            <span class="text-sm text-gray-700">Enriched only</span>
          </label>
        </div>
      </div>
    </div>
    
    <div class="flex justify-between items-center mt-6">
      <button
        @click="clearFilters"
        class="text-sm text-gray-600 hover:text-gray-800"
      >
        Clear all filters
      </button>
      
      <div class="flex space-x-3">
        <button
          @click="resetFilters"
          class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50"
        >
          Reset
        </button>
        <button
          @click="applyFilters"
          class="px-4 py-2 border border-transparent rounded-md text-sm font-medium text-white bg-blue-600 hover:bg-blue-700"
        >
          Apply Filters
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { AppFilters, Stock } from '@/types/api'

interface Props {
  filters: AppFilters
  stocks: Stock[]
}

interface Emits {
  (e: 'update:filters', filters: AppFilters): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const localFilters = ref<AppFilters>({ ...props.filters })

const availableSectors = computed(() => {
  const sectors = new Set<string>()
  props.stocks.forEach(stock => {
    if (stock.sector) sectors.add(stock.sector)
  })
  return Array.from(sectors).sort()
})

const availableActions = computed(() => {
  const actions = new Set<string>()
  props.stocks.forEach(stock => {
    if (stock.action) actions.add(stock.action)
  })
  return Array.from(actions).sort()
})

const availableBrokerages = computed(() => {
  const brokerages = new Set<string>()
  props.stocks.forEach(stock => {
    if (stock.brokerage) brokerages.add(stock.brokerage)
  })
  return Array.from(brokerages).sort()
})

const clearFilters = () => {
  localFilters.value = {
    sector: [],
    recommendation_type: [],
    action: [],
    brokerage: [],
    enriched_only: false
  }
}

const resetFilters = () => {
  localFilters.value = { ...props.filters }
}

const applyFilters = () => {
  emit('update:filters', { ...localFilters.value })
}

watch(() => props.filters, (newFilters) => {
  localFilters.value = { ...newFilters }
}, { deep: true })
</script>