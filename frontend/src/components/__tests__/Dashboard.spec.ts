import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount } from '@vue/test-utils'
import { createPinia, setActivePinia } from 'pinia'
import Dashboard from '../Dashboard.vue'
import DashboardCard from '../DashboardCard.vue'
import { useStocksStore } from '@/stores/stocks'

// Mock DashboardCard component
vi.mock('../DashboardCard.vue', () => ({
  default: {
    name: 'DashboardCard',
    props: ['title', 'value', 'icon', 'subtitle'],
    template: '<div class="dashboard-card" :data-testid="`card-${title?.toLowerCase().replace(/\\s+/g, \'-\')}`">{{ title }}: {{ value }}</div>'
  }
}))

// Mock the API service
vi.mock('@/services/api', () => ({
  apiService: {
    getStocks: vi.fn(),
    getRecommendations: vi.fn(),
    enrichStocks: vi.fn(),
    syncStocks: vi.fn(),
    checkHealth: vi.fn()
  }
}))

describe('Dashboard', () => {
  let wrapper: any
  let store: any

  beforeEach(() => {
    setActivePinia(createPinia())
    store = useStocksStore()
    
    // Mock store data
    store.stocks = [
      { 
        id: 1, 
        ticker: 'AAPL', 
        company: 'Apple Inc.', 
        sector: 'Technology',
        enriched_at: '2024-01-01T00:00:00Z'
      },
      { 
        id: 2, 
        ticker: 'GOOGL', 
        company: 'Google LLC', 
        sector: 'Technology',
        enriched_at: null
      },
      { 
        id: 3, 
        ticker: 'JPM', 
        company: 'JPMorgan Chase', 
        sector: 'Financial Services',
        enriched_at: '2024-01-01T00:00:00Z'
      }
    ]
    
    store.recommendations = [
      {
        symbol: 'AAPL',
        name: 'Apple Inc.',
        score: 85.0,
        reason: 'Strong buy signal',
        recommendation_type: 'BUY'
      },
      {
        symbol: 'GOOGL',
        name: 'Google LLC',
        score: 70.0,
        reason: 'Hold recommendation',
        recommendation_type: 'HOLD'
      }
    ]

    wrapper = mount(Dashboard)
  })

  it('renders dashboard title', () => {
    expect(wrapper.find('h1').text()).toBe('Dashboard')
  })

  it('displays dashboard cards with correct data', async () => {
    await wrapper.vm.$nextTick()
    
    const cards = wrapper.findAllComponents({ name: 'DashboardCard' })
    expect(cards).toHaveLength(4)
    
    // Check Total Stocks card
    const totalStocksCard = wrapper.find('[data-testid="card-total-stocks"]')
    expect(totalStocksCard.exists()).toBe(true)
    expect(totalStocksCard.text()).toContain('Total Stocks: 3')
    
    // Check Enriched Data card
    const enrichedCard = wrapper.find('[data-testid="card-enriched-data"]')
    expect(enrichedCard.exists()).toBe(true)
    expect(enrichedCard.text()).toContain('Enriched Data: 2')
  })

  it('calculates enrichment percentage correctly', () => {
    expect(wrapper.vm.enrichmentPercentage).toBe(67) // 2/3 * 100 = 67%
  })

  it('calculates buy recommendations count correctly', () => {
    expect(wrapper.vm.buyRecommendationsCount).toBe(1)
  })

  it('displays sector breakdown correctly', async () => {
    await wrapper.vm.$nextTick()
    
    // Check if the component text contains the section title and data
    expect(wrapper.text()).toContain('Sector Breakdown')
    expect(wrapper.text()).toContain('Technology')
    expect(wrapper.text()).toContain('Financial Services')
  })

  it('displays recommendations correctly', async () => {
    await wrapper.vm.$nextTick()
    
    // Check if the component text contains the section title and data
    expect(wrapper.text()).toContain('Top Recommendations')
    expect(wrapper.text()).toContain('AAPL')
    expect(wrapper.text()).toContain('GOOGL')
    expect(wrapper.text()).toContain('BUY')
    expect(wrapper.text()).toContain('HOLD')
  })

  it('returns correct recommendation badge colors', () => {
    expect(wrapper.vm.getRecommendationBadgeColor('BUY')).toBe('bg-green-100 text-green-800')
    expect(wrapper.vm.getRecommendationBadgeColor('HOLD')).toBe('bg-yellow-100 text-yellow-800')
    expect(wrapper.vm.getRecommendationBadgeColor('WATCH')).toBe('bg-blue-100 text-blue-800')
    expect(wrapper.vm.getRecommendationBadgeColor('SELL')).toBe('bg-red-100 text-red-800')
    expect(wrapper.vm.getRecommendationBadgeColor('UNKNOWN')).toBe('bg-gray-100 text-gray-800')
  })

  it('calls refresh data on refresh button click', async () => {
    const fetchStocksSpy = vi.spyOn(store, 'fetchStocks')
    const fetchRecommendationsSpy = vi.spyOn(store, 'fetchRecommendations')
    
    const refreshButton = wrapper.findAll('button').find((w: any) => w.text().includes('Refresh'))
    await refreshButton?.trigger('click')
    
    expect(fetchStocksSpy).toHaveBeenCalled()
    expect(fetchRecommendationsSpy).toHaveBeenCalledWith(10)
  })

  it('calls enrich data on enrich button click', async () => {
    const enrichStocksSpy = vi.spyOn(store, 'enrichStocks')
    
    const enrichButton = wrapper.findAll('button').find((w: any) => w.text().includes('Enrich Data'))
    await enrichButton?.trigger('click')
    
    expect(enrichStocksSpy).toHaveBeenCalledWith(10)
  })

  it('disables buttons when loading', async () => {
    store.loading.stocks = true
    store.loading.recommendations = true
    await wrapper.vm.$nextTick()
    
    const refreshButton = wrapper.findAll('button').find((w: any) => w.text().includes('Refresh'))
    expect(refreshButton?.attributes('disabled')).toBeDefined()
  })

  it('shows error message when error exists', async () => {
    store.error = 'Test error message'
    await wrapper.vm.$nextTick()
    
    const errorDiv = wrapper.find('.bg-red-50\\/80')
    expect(errorDiv.exists()).toBe(true)
    expect(errorDiv.text()).toContain('Test error message')
  })

  it('shows enrichment stats when available', async () => {
    store.enrichmentStats = {
      total_stocks: 10,
      newly_enriched: 5,
      already_enriched: 3,
      failed: 2
    }
    await wrapper.vm.$nextTick()
    
    expect(wrapper.text()).toContain('Latest Enrichment Results')
    expect(wrapper.text()).toContain('10') // total stocks
    expect(wrapper.text()).toContain('5')  // newly enriched
    expect(wrapper.text()).toContain('3')  // already enriched
    expect(wrapper.text()).toContain('2')  // failed
  })

  it('shows no data message when no sectors available', async () => {
    store.stocks = []
    await wrapper.vm.$nextTick()
    
    expect(wrapper.text()).toContain('No sector data available')
  })

  it('shows no data message when no recommendations available', async () => {
    store.recommendations = []
    await wrapper.vm.$nextTick()
    
    expect(wrapper.text()).toContain('No recommendations available')
  })

  it('limits sector breakdown to 8 items', async () => {
    // Create 10 different sectors
    store.stocks = Array.from({ length: 10 }, (_, i) => ({
      id: i + 1,
      ticker: `STOCK${i}`,
      company: `Company ${i}`,
      sector: `Sector ${i}`,
      enriched_at: null
    }))
    
    await wrapper.vm.$nextTick()
    
    // Check the computed sectorBreakdown is limited in the template
    // The template uses .slice(0, 8) so we test the rendered sector bars
    const sectorBars = wrapper.findAll('.w-full.bg-slate-100.rounded-full.h-2\\.5')
    expect(sectorBars.length).toBeLessThanOrEqual(8)
  })

  it('limits recommendations to 5 items', async () => {
    // Create 7 recommendations
    store.recommendations = Array.from({ length: 7 }, (_, i) => ({
      symbol: `STOCK${i}`,
      name: `Company ${i}`,
      score: 80.0,
      reason: 'Test reason',
      recommendation_type: 'BUY'
    }))
    
    await wrapper.vm.$nextTick()
    
    // Should only show 5 recommendations max
    const recommendationElements = wrapper.findAll('.p-5.rounded-xl.border.border-slate-200')
    expect(recommendationElements.length).toBeLessThanOrEqual(5)
  })
})