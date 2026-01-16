package cch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ModelStat represents statistics for a single model
type ModelStat struct {
	Model     string  `json:"model"`
	TotalCost float64 `json:"totalCost"`
}

// KeyInfo represents a single API key's information
type KeyInfo struct {
	FullKey               string      `json:"fullKey"`
	TodayUsage            float64     `json:"todayUsage"`
	TodayCallCount        int         `json:"todayCallCount"`
	Limit5hUsd            float64     `json:"limit5hUsd"`
	LimitWeeklyUsd        float64     `json:"limitWeeklyUsd"`
	LimitMonthlyUsd       float64     `json:"limitMonthlyUsd"`
	LimitConcurrentSessions int       `json:"limitConcurrentSessions"`
	LastProviderName      string      `json:"lastProviderName"`
	ModelStats            []ModelStat `json:"modelStats"`
}

// UserData represents a user's data from the API
type UserData struct {
	DailyQuota float64   `json:"dailyQuota"`
	Keys       []KeyInfo `json:"keys"`
}

// GetUsersResponse represents the API response for getUsers
type GetUsersResponse struct {
	OK    bool       `json:"ok"`
	Data  []UserData `json:"data"`
	Error string     `json:"error,omitempty"`
}

// Stats represents the collected statistics
type Stats struct {
	TodayCost               float64
	TodayRequests           int
	DailyQuota              float64
	Limit5h                 float64
	LimitWeekly             float64
	LimitMonthly            float64
	LimitConcurrentSessions int
	LastProviderName        string
	LastUsedModel           string
}

// cacheEntry holds cached data with expiry
type cacheEntry struct {
	data   *Stats
	expiry time.Time
}

// Client is the CCH API client
type Client struct {
	baseURL  string
	apiKey   string
	client   *http.Client
	cache    *cacheEntry
	cacheTTL time.Duration
	mu       sync.RWMutex
}

// NewClient creates a new CCH client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:  baseURL,
		apiKey:   apiKey,
		client:   &http.Client{Timeout: 10 * time.Second},
		cacheTTL: 7 * time.Second,
	}
}

// GetStats fetches statistics from the CCH API with caching
func (c *Client) GetStats() (*Stats, error) {
	// Check cache first
	c.mu.RLock()
	if c.cache != nil && time.Now().Before(c.cache.expiry) {
		stats := c.cache.data
		c.mu.RUnlock()
		return stats, nil
	}
	c.mu.RUnlock()

	// Fetch fresh data
	stats, err := c.fetchStats()
	if err != nil {
		return nil, err
	}

	// Update cache
	c.mu.Lock()
	c.cache = &cacheEntry{
		data:   stats,
		expiry: time.Now().Add(c.cacheTTL),
	}
	c.mu.Unlock()

	return stats, nil
}

// fetchStats fetches statistics from the API
func (c *Client) fetchStats() (*Stats, error) {
	// Get user data (auth-token cookie is simply the API key)
	userData, err := c.getUserData()
	if err != nil {
		return nil, err
	}

	if len(userData) == 0 {
		return nil, fmt.Errorf("no user data available")
	}

	user := userData[0]

	// Find our key by matching the API key
	var ourKey *KeyInfo
	for i := range user.Keys {
		if user.Keys[i].FullKey == c.apiKey {
			ourKey = &user.Keys[i]
			break
		}
	}

	if ourKey == nil {
		return nil, fmt.Errorf("key not found in user data")
	}

	// Extract the model with highest cost (likely the most recently used expensive model)
	var lastUsedModel string
	if len(ourKey.ModelStats) > 0 {
		maxCost := ourKey.ModelStats[0].TotalCost
		lastUsedModel = ourKey.ModelStats[0].Model
		for _, ms := range ourKey.ModelStats[1:] {
			if ms.TotalCost > maxCost {
				maxCost = ms.TotalCost
				lastUsedModel = ms.Model
			}
		}
	}

	return &Stats{
		TodayCost:               ourKey.TodayUsage,
		TodayRequests:           ourKey.TodayCallCount,
		DailyQuota:              user.DailyQuota,
		Limit5h:                 ourKey.Limit5hUsd,
		LimitWeekly:             ourKey.LimitWeeklyUsd,
		LimitMonthly:            ourKey.LimitMonthlyUsd,
		LimitConcurrentSessions: ourKey.LimitConcurrentSessions,
		LastProviderName:        ourKey.LastProviderName,
		LastUsedModel:           lastUsedModel,
	}, nil
}

// getUserData fetches user data from the API
func (c *Client) getUserData() ([]UserData, error) {
	url := fmt.Sprintf("%s/api/actions/users/getUsers", c.baseURL)

	reqBody, _ := json.Marshal(map[string]any{})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("auth-token=%s", c.apiKey))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var result GetUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, fmt.Errorf("API error: %s", result.Error)
	}

	return result.Data, nil
}

// ClearCache clears the cached data
func (c *Client) ClearCache() {
	c.mu.Lock()
	c.cache = nil
	c.mu.Unlock()
}
