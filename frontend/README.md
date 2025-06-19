# Stock Recommender Frontend

A Vue 3 + TypeScript frontend application for the Stock Recommender API, featuring intelligent stock analysis and recommendations.

## Features

### 🎯 Dashboard
- **Metrics Cards**: Total stocks, enriched data percentage, top recommendations, sectors tracked
- **Sector Breakdown**: Visual representation of stock distribution across sectors
- **Top Recommendations**: Quick view of highest-scoring investment opportunities
- **Enrichment Progress**: Real-time stats from data enrichment operations

### 📊 Stock Management
- **Comprehensive Stock List**: Paginated table with filtering and sorting
- **Advanced Filtering**: Filter by sector, action, brokerage, enrichment status
- **Real-time Data**: Sync and enrich stock data with Alpha Vantage API
- **Market Data Display**: Market cap, P/E ratios, 52-week ranges

### 🎯 Smart Recommendations
- **AI-Powered Scoring**: Machine learning-based recommendation scores (0-100)
- **Visual Score Indicators**: Color-coded progress bars and badges
- **Upside Calculation**: Automatic target vs current price analysis
- **Filtering & Sorting**: Multiple criteria for finding the best opportunities

### 🎨 UI/UX Features
- **Responsive Design**: Mobile-first approach with Tailwind CSS
- **Real-time Loading States**: Smooth user experience with loading indicators
- **Error Handling**: Comprehensive error display and retry mechanisms
- **Modern Design**: Clean, professional interface with intuitive navigation

## Tech Stack

- **Framework**: Vue 3 with Composition API
- **Language**: TypeScript for type safety
- **State Management**: Pinia for reactive global state
- **Styling**: Tailwind CSS for modern, responsive design
- **Build Tool**: Vite for fast development and optimized builds
- **Testing**: Vitest for unit testing

## API Integration

### Endpoints Used
- `GET /health` - Service health check
- `GET /api/stocks` - Fetch all stocks with pagination
- `POST /api/stocks/sync` - Synchronize stock data
- `POST /api/stocks/enrich` - Enrich stocks with fundamental data
- `GET /api/recommendations` - Get AI-powered recommendations

### Data Flow
1. **Initialization**: Health check → Load stocks → Load recommendations
2. **Data Enrichment**: Manual trigger → API call → Progress tracking → Refresh
3. **Filtering**: Client-side filtering for instant responsiveness
4. **Real-time Updates**: Automatic refresh after sync/enrichment operations

## Project Setup

```sh
npm install
```

### Compile and Hot-Reload for Development

```sh
npm run dev
```

### Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Run Unit Tests with [Vitest](https://vitest.dev/)

```sh
npm run test:unit
```

## Environment Variables
- `VITE_API_URL`: Backend API URL (default: http://localhost:8080)

## Project Structure

```
src/
├── components/
│   ├── Dashboard.vue              # Main dashboard with metrics
│   ├── DashboardCard.vue          # Reusable metric card component
│   ├── StocksList.vue             # Stock table with filtering
│   ├── StockFilters.vue           # Advanced filtering controls
│   ├── StockCard.vue              # Individual stock display card
│   ├── Recommendations.vue        # Smart recommendations view
│   └── RecommendationCard.vue     # Individual recommendation card
├── views/
│   ├── HomeView.vue              # Dashboard page
│   ├── StocksView.vue            # Stocks management page
│   └── RecommendationsView.vue   # Recommendations page
├── stores/
│   └── stocks.ts                 # Pinia store for state management
├── services/
│   └── api.ts                    # API service layer
├── composables/
│   └── useApi.ts                 # Reusable API composables
├── types/
│   └── api.ts                    # TypeScript interfaces
└── router/
    └── index.ts                  # Vue Router configuration
```

## Key Features Implementation

### State Management (Pinia)
- Centralized state for stocks, recommendations, loading states
- Computed properties for filtering and pagination
- Async actions for API calls with error handling

### Composables Pattern
- Reusable API logic with `useApi` composables
- Consistent error handling and loading states
- Type-safe API interactions

### Component Architecture
- Single File Components with `<script setup>` syntax
- Props and emits with TypeScript interfaces
- Composition API for reactive logic

## Color Scheme & Design

```css
/* Recommendation Types */
BUY:   bg-green-100 text-green-800
HOLD:  bg-yellow-100 text-yellow-800  
WATCH: bg-blue-100 text-blue-800
SELL:  bg-red-100 text-red-800

/* Score Indicators */
80-100: bg-green-500 (Excellent)
60-79:  bg-yellow-500 (Good)
40-59:  bg-orange-500 (Fair)
0-39:   bg-red-500 (Poor)
```

## Browser Support

- Chrome 88+
- Firefox 85+
- Safari 14+
- Edge 88+

## Recommended IDE Setup

[VSCode](https://code.visualstudio.com/) + [Volar](https://marketplace.visualstudio.com/items?itemName=Vue.volar) (and disable Vetur).
