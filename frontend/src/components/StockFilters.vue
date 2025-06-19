<template>
  <div class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-6 shadow-sm">
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-6 gap-4 items-end">
      <!-- Search Input -->
      <div class="lg:col-span-2">
        <label for="search" class="block text-sm font-medium text-slate-700 mb-2">Search</label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="h-4 w-4 text-slate-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <input
            id="search"
            v-model="localFilters.search"
            type="text"
            placeholder="Search ticker or company..."
            class="w-full pl-10 pr-4 py-3 bg-white border border-slate-200 rounded-xl text-sm placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200"
            @input="onFilterChange"
          />
        </div>
      </div>

      <!-- Sector Filter -->
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">Sector</label>
        <div class="relative">
          <select 
            v-model="localFilters.sector" 
            @change="onFilterChange" 
            class="w-full appearance-none bg-white border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 cursor-pointer"
          >
            <option value="" class="text-slate-500">All Sectors</option>
            <option v-for="sector in availableSectors" :key="sector" :value="sector" class="text-slate-700">{{ sector }}</option>
          </select>
          <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
            <svg class="h-4 w-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Action Filter -->
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">Action</label>
        <div class="relative">
          <select 
            v-model="localFilters.action" 
            @change="onFilterChange" 
            class="w-full appearance-none bg-white border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 cursor-pointer"
          >
            <option value="" class="text-slate-500">All Actions</option>
            <option v-for="action in availableActions" :key="action" :value="action" class="text-slate-700">{{ action }}</option>
          </select>
          <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
            <svg class="h-4 w-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Brokerage Filter -->
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">Brokerage</label>
        <div class="relative">
          <select 
            v-model="localFilters.brokerage" 
            @change="onFilterChange" 
            class="w-full appearance-none bg-white border border-slate-200 rounded-xl px-4 py-3 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-blue-500/20 focus:border-blue-500 transition-all duration-200 cursor-pointer"
          >
            <option value="" class="text-slate-500">All Brokerages</option>
            <option v-for="brokerage in availableBrokerages" :key="brokerage" :value="brokerage" class="text-slate-700">{{ brokerage }}</option>
          </select>
          <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
            <svg class="h-4 w-4 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Controls -->
      <div class="flex flex-col space-y-3">
        <label class="flex items-center group cursor-pointer">
          <div class="relative">
            <input 
              type="checkbox" 
              v-model="localFilters.enriched_only" 
              @change="onFilterChange" 
              class="sr-only peer"
            >
            <div class="w-11 h-6 bg-slate-200 peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-500/20 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-500"></div>
          </div>
          <span class="ml-3 text-sm font-medium text-slate-700">Enriched Only</span>
        </label>
        <button 
          @click="clearFilters" 
          class="inline-flex items-center justify-center px-4 py-2 text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-50 rounded-xl transition-all duration-200"
        >
          <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
          Clear
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import type { AppFilters, Stock } from '@/types/api'
import { debounce } from 'lodash-es'

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

const onFilterChange = debounce(() => {
  emit('update:filters', { ...localFilters.value })
}, 300)

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
    search: '',
    sector: '',
    recommendation_type: '',
    action: '',
    brokerage: '',
    enriched_only: false
  }
  onFilterChange()
}

watch(() => props.filters, (newFilters) => {
  localFilters.value = { ...newFilters }
}, { deep: true })
</script>