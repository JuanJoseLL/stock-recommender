package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	AlphaVantageBaseURL = "https://www.alphavantage.co/query"
)

type AlphaVantageClient struct {
	apiKey     string
	httpClient *http.Client
}

type AlphaVantageTimeSeriesResponse struct {
	MetaData   AlphaVantageMetaData            `json:"Meta Data"`
	TimeSeries map[string]AlphaVantageDataPoint `json:"Time Series (Daily)"`
}

type AlphaVantageMetaData struct {
	Information   string `json:"1. Information"`
	Symbol        string `json:"2. Symbol"`
	LastRefreshed string `json:"3. Last Refreshed"`
	OutputSize    string `json:"4. Output Size"`
	TimeZone      string `json:"5. Time Zone"`
}

type AlphaVantageDataPoint struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type AlphaVantageCompanyOverview struct {
	Symbol             string `json:"Symbol"`
	AssetType          string `json:"AssetType"`
	Name               string `json:"Name"`
	Description        string `json:"Description"`
	CIK                string `json:"CIK"`
	Exchange           string `json:"Exchange"`
	Currency           string `json:"Currency"`
	Country            string `json:"Country"`
	Sector             string `json:"Sector"`
	Industry           string `json:"Industry"`
	Address            string `json:"Address"`
	FullTimeEmployees  string `json:"FullTimeEmployees"`
	FiscalYearEnd      string `json:"FiscalYearEnd"`
	LatestQuarter      string `json:"LatestQuarter"`
	MarketCapitalization string `json:"MarketCapitalization"`
	EBITDA             string `json:"EBITDA"`
	PERatio            string `json:"PERatio"`
	PEGRatio           string `json:"PEGRatio"`
	BookValue          string `json:"BookValue"`
	DividendPerShare   string `json:"DividendPerShare"`
	DividendYield      string `json:"DividendYield"`
	EPS                string `json:"EPS"`
	RevenuePerShareTTM string `json:"RevenuePerShareTTM"`
	ProfitMargin       string `json:"ProfitMargin"`
	OperatingMarginTTM string `json:"OperatingMarginTTM"`
	ReturnOnAssetsTTM  string `json:"ReturnOnAssetsTTM"`
	ReturnOnEquityTTM  string `json:"ReturnOnEquityTTM"`
	RevenueTTM         string `json:"RevenueTTM"`
	GrossProfitTTM     string `json:"GrossProfitTTM"`
	DilutedEPSTTM      string `json:"DilutedEPSTTM"`
	QuarterlyEarningsGrowthYOY string `json:"QuarterlyEarningsGrowthYOY"`
	QuarterlyRevenueGrowthYOY  string `json:"QuarterlyRevenueGrowthYOY"`
	AnalystTargetPrice         string `json:"AnalystTargetPrice"`
	TrailingPE                 string `json:"TrailingPE"`
	ForwardPE                  string `json:"ForwardPE"`
	PriceToSalesRatioTTM       string `json:"PriceToSalesRatioTTM"`
	PriceToBookRatio           string `json:"PriceToBookRatio"`
	EVToRevenue                string `json:"EVToRevenue"`
	EVToEBITDA                 string `json:"EVToEBITDA"`
	Beta                       string `json:"Beta"`
	WeekHigh52                 string `json:"52WeekHigh"`
	WeekLow52                  string `json:"52WeekLow"`
	MovingAverage50            string `json:"50DayMovingAverage"`
	MovingAverage200           string `json:"200DayMovingAverage"`
	SharesOutstanding          string `json:"SharesOutstanding"`
	DividendDate               string `json:"DividendDate"`
	ExDividendDate             string `json:"ExDividendDate"`
}

type AlphaVantageTopGainersLosers struct {
	MetaData                string                        `json:"metadata"`
	LastUpdated             string                        `json:"last_updated"`
	TopGainers              []AlphaVantageMarketMover     `json:"top_gainers"`
	TopLosers               []AlphaVantageMarketMover     `json:"top_losers"`
	MostActivelyTraded      []AlphaVantageMarketMover     `json:"most_actively_traded"`
}

type AlphaVantageMarketMover struct {
	Ticker           string `json:"ticker"`
	Price            string `json:"price"`
	ChangeAmount     string `json:"change_amount"`
	ChangePercentage string `json:"change_percentage"`
	Volume           string `json:"volume"`
}

func NewAlphaVantageClient(apiKey string) *AlphaVantageClient {
	return &AlphaVantageClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *AlphaVantageClient) GetTimeSeries(ctx context.Context, symbol string) (*AlphaVantageTimeSeriesResponse, error) {
	url := fmt.Sprintf("%s?function=TIME_SERIES_DAILY&symbol=%s&apikey=%s", 
		AlphaVantageBaseURL, symbol, c.apiKey)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var timeSeriesResp AlphaVantageTimeSeriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&timeSeriesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &timeSeriesResp, nil
}

func (c *AlphaVantageClient) GetCompanyOverview(ctx context.Context, symbol string) (*AlphaVantageCompanyOverview, error) {
	url := fmt.Sprintf("%s?function=OVERVIEW&symbol=%s&apikey=%s", 
		AlphaVantageBaseURL, symbol, c.apiKey)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var overview AlphaVantageCompanyOverview
	if err := json.NewDecoder(resp.Body).Decode(&overview); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &overview, nil
}

func (c *AlphaVantageClient) GetTopGainersLosers(ctx context.Context) (*AlphaVantageTopGainersLosers, error) {
	url := fmt.Sprintf("%s?function=TOP_GAINERS_LOSERS&apikey=%s", 
		AlphaVantageBaseURL, c.apiKey)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var topGainersLosers AlphaVantageTopGainersLosers
	if err := json.NewDecoder(resp.Body).Decode(&topGainersLosers); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &topGainersLosers, nil
}