<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="min-h-screen bg-gray-50">
    <nav class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between h-16">
          <div class="flex items-center">
            <div class="flex-shrink-0 flex items-center">
              <h1 class="text-xl font-bold text-gray-900">Stock Recommender</h1>
            </div>
            <div class="hidden sm:ml-8 sm:flex sm:space-x-8">
              <RouterLink
                to="/"
                class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors"
                active-class="border-blue-500 text-blue-600"
              >
                Dashboard
              </RouterLink>
              <RouterLink
                to="/stocks"
                class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors"
                active-class="border-blue-500 text-blue-600"
              >
                Stocks
              </RouterLink>
              <RouterLink
                to="/recommendations"
                class="border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors"
                active-class="border-blue-500 text-blue-600"
              >
                Recommendations
              </RouterLink>
            </div>
          </div>
          <div class="flex items-center">
            <div class="text-sm text-gray-500">
              <span class="hidden sm:inline">AI-Powered Stock Analysis</span>
            </div>
          </div>
        </div>
      </div>
    </nav>
    
    <RouterView />
  </div>
</template>

<style scoped>
.router-link-active {
  border-color: #3b82f6;
  color: #2563eb;
}
</style>
