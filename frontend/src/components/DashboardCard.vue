<template>
  <div class="bg-white/80 backdrop-blur-sm rounded-2xl border border-slate-200/50 p-6 shadow-sm">
    <div class="flex items-center space-x-4">
      <div class="w-14 h-14 rounded-2xl flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-50 text-blue-600 flex-shrink-0">
        <svg class="w-7 h-7" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path v-if="icon === 'chart'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
          <path v-else-if="icon === 'check'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
          <path v-else-if="icon === 'target'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"></path>
          <path v-else-if="icon === 'building'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h2M9 21V9m6 12V9"></path>
          <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
        </svg>
      </div>
      
      <div class="flex-1 min-w-0">
        <p class="text-sm font-medium text-slate-600 mb-1">{{ title }}</p>
        <p class="text-3xl font-bold text-slate-900 leading-tight">{{ value }}</p>
        <div v-if="subtitle" class="flex items-center space-x-2 mt-2">
          <div class="w-2 h-2 rounded-full bg-emerald-400"></div>
          <p class="text-sm text-slate-500 truncate">{{ subtitle }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title: string
  value: string | number
  subtitle?: string
  icon?: string
  trend?: {
    value: string
    label: string
    positive: boolean
  }
}

defineProps<Props>()
</script>