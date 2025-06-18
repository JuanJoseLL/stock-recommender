package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/JuanJoseLL/stock-recommender/internal/domain"
)

// ExternalAPIClient defines the interface for fetching stock data from external APIs
type ExternalAPIClient interface {
	FetchStocks(ctx context.Context) (*domain.APIResponse, error)
	FetchStocksWithPagination(ctx context.Context, nextPage string) (*domain.APIResponse, error)
	FetchAllStocks(ctx context.Context) ([]domain.Stock, error)
}

// HTTPClient implements ExternalAPIClient for HTTP-based stock APIs
type HTTPClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewHTTPClient creates a new HTTP client for external stock APIs
func NewHTTPClient(baseURL, apiKey string) *HTTPClient {
	return &HTTPClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchStocks retrieves stock data from the external API (first page)
func (c *HTTPClient) FetchStocks(ctx context.Context) (*domain.APIResponse, error) {
	return c.FetchStocksWithPagination(ctx, "")
}

// FetchStocksWithPagination retrieves stock data with pagination support
func (c *HTTPClient) FetchStocksWithPagination(ctx context.Context, nextPage string) (*domain.APIResponse, error) {
	reqURL := c.baseURL + "/production/swechallenge/list"
	
	// Add pagination parameter if provided
	if nextPage != "" {
		parsedURL, err := url.Parse(reqURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}
		params := parsedURL.Query()
		params.Set("next_page", nextPage)
		parsedURL.RawQuery = params.Encode()
		reqURL = parsedURL.String()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResponse domain.APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse, nil
}

// FetchAllStocks retrieves all stock data by following pagination
func (c *HTTPClient) FetchAllStocks(ctx context.Context) ([]domain.Stock, error) {
	var allStocks []domain.Stock
	nextPage := ""
	page := 1
	const maxPages = 1000 // Safety limit to prevent infinite loops

	for page <= maxPages {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		response, err := c.FetchStocksWithPagination(ctx, nextPage)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page %d: %w", page, err)
		}

		if len(response.Items) == 0 {
			break // No more data
		}

		allStocks = append(allStocks, response.Items...)

		// Check if there's a next page
		if response.NextPage == "" {
			break // No more pages
		}

		nextPage = response.NextPage
		page++

		// Add a small delay to be respectful to the API
		time.Sleep(100 * time.Millisecond)
	}

	if page > maxPages {
		return nil, fmt.Errorf("reached maximum page limit (%d), possible infinite pagination", maxPages)
	}

	return allStocks, nil
}