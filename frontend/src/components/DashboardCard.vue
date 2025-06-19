<template>
  <div class="bg-white rounded-lg shadow-md p-6 border border-gray-200">
    <div class="flex items-center justify-between">
      <div>
        <p class="text-sm font-medium text-gray-600">{{ title }}</p>
        <p class="text-3xl font-bold text-gray-900 mt-2">{{ value }}</p>
        <p v-if="subtitle" class="text-sm text-gray-500 mt-1">{{ subtitle }}</p>
      </div>
      <div v-if="icon" class="flex-shrink-0">
        <div class="w-12 h-12 rounded-lg flex items-center justify-center" :class="iconBgColor">
          <svg class="w-6 h-6" :class="iconColor" fill="currentColor" viewBox="0 0 24 24">
            <path v-if="icon === 'chart'" d="M3 3v18h18v-2H5V3H3zm3 14h2V9H6v8zm4 0h2V5h-2v12zm4 0h2v-6h-2v6z"/>
            <path v-else-if="icon === 'check'" d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
            <path v-else-if="icon === 'target'" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
            <path v-else-if="icon === 'building'" d="M6.5 14.5v-3.5c0-.83.67-1.5 1.5-1.5s1.5.67 1.5 1.5v3.5h1v-3.5c0-.83.67-1.5 1.5-1.5s1.5.67 1.5 1.5v3.5h1v-3.5c0-1.38-1.12-2.5-2.5-2.5-.52 0-1 .16-1.39.43-.39-.27-.87-.43-1.39-.43C6.12 9 5 10.12 5 11.5v3.5h1.5z"/>
            <path v-else d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z"/>
          </svg>
        </div>
      </div>
    </div>
    <div v-if="trend" class="mt-4 flex items-center">
      <span 
        class="inline-flex items-center text-sm font-medium"
        :class="trend.positive ? 'text-green-600' : 'text-red-600'"
      >
        <span class="mr-1">{{ trend.positive ? '↗' : '↘' }}</span>
        {{ trend.value }}
      </span>
      <span class="text-gray-500 text-sm ml-2">{{ trend.label }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title: string
  value: string | number
  subtitle?: string
  icon?: string
  iconColor?: string
  iconBgColor?: string
  trend?: {
    value: string
    label: string
    positive: boolean
  }
}

withDefaults(defineProps<Props>(), {
  iconColor: 'text-blue-600',
  iconBgColor: 'bg-blue-100'
})
</script>