package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/JuanJoseLL/stock-recommender/internal/domain"
)

// ExternalAPIClient defines the interface for fetching stock data from external APIs
type ExternalAPIClient interface {
	FetchStocks(ctx context.Context) (*domain.APIResponse, error)
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

// FetchStocks retrieves stock data from the external API
func (c *HTTPClient) FetchStocks(ctx context.Context) (*domain.APIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/production/swechallenge/list", nil)
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