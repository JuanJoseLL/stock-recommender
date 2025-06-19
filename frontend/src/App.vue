<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <div class="flex h-screen bg-gradient-to-br from-slate-50 to-slate-100 font-sans">
    <!-- Sidebar -->
    <aside class="w-64 bg-white/80 backdrop-blur-sm border-r border-slate-200/50 flex-shrink-0">
      <div class="p-6 border-b border-slate-200/50">
        <h1 class="text-xl font-bold text-slate-900">Stock Recommender</h1>
        <span class="text-sm text-slate-600">AI-Powered Analysis</span>
      </div>
      <nav class="mt-6 px-3">
        <RouterLink
          to="/"
          class="flex items-center py-3 px-4 mb-2 text-slate-700 hover:bg-slate-100/70 hover:text-slate-900 rounded-xl transition-all duration-200 group"
          active-class="bg-slate-900 text-white shadow-sm"
        >
          <svg class="w-5 h-5 mr-3 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10v11h6v-7h6v7h6V10L12 3z"></path>
          </svg>
          <span class="font-medium">Dashboard</span>
        </RouterLink>
        <RouterLink
          to="/stocks"
          class="flex items-center py-3 px-4 mb-2 text-slate-700 hover:bg-slate-100/70 hover:text-slate-900 rounded-xl transition-all duration-200 group"
          active-class="bg-slate-900 text-white shadow-sm"
        >
          <svg class="w-5 h-5 mr-3 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
          </svg>
          <span class="font-medium">Stocks</span>
        </RouterLink>
        <RouterLink
          to="/recommendations"
          class="flex items-center py-3 px-4 mb-2 text-slate-700 hover:bg-slate-100/70 hover:text-slate-900 rounded-xl transition-all duration-200 group"
          active-class="bg-slate-900 text-white shadow-sm"
        >
          <svg class="w-5 h-5 mr-3 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 12l3-3 3 3 4-4M8 21l4-4 4 4M3 4h18M4 4h16v12a1 1 0 01-1 1H5a1 1 0 01-1-1V4z"></path>
          </svg>
          <span class="font-medium">Recommendations</span>
        </RouterLink>
      </nav>
    </aside>

    <!-- Main content -->
    <main class="flex-1 flex flex-col overflow-hidden">
      <div class="flex-1 overflow-x-hidden overflow-y-auto">
        <div class="container mx-auto px-6 py-8">
          <RouterView />
        </div>
      </div>
    </main>
  </div>
</template>

<style>
/* Modern focus styles */
:focus {
  outline: none;
}

/* Custom scrollbar for modern look */
::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}
::-webkit-scrollbar-track {
  background: transparent;
}
::-webkit-scrollbar-thumb {
  background: rgb(148 163 184 / 0.3);
  border-radius: 3px;
}
::-webkit-scrollbar-thumb:hover {
  background: rgb(148 163 184 / 0.5);
}

/* Active link styling */
.router-link-exact-active {
  background-color: rgb(15 23 42) !important;
  color: #ffffff !important;
}

/* Line clamp utilities for text truncation */
.line-clamp-1 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 1;
}

.line-clamp-2 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
}

.line-clamp-3 {
  overflow: hidden;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
}
</style>
