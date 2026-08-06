package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Profile struct {
	UserID     uint64 `json:"user_id"`
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}

type Metrics struct {
	ViewsTotal        int `json:"viewsTotal"`
	FavoritesReceived int `json:"favoritesReceived"`
	DealsClosed       int `json:"dealsClosed"`
	DaysActive        int `json:"daysActive"`
}

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *UserClient) GetProfile(ctx context.Context, id string) (Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/%s", c.baseURL, id), nil)
	if err != nil {
		return Profile{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return Profile{}, fmt.Errorf("profile not found")
	}
	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var profile Profile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return Profile{}, err
	}

	return profile, nil
}

type AnalyticsClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewAnalyticsClient(baseURL string) *AnalyticsClient {
	return &AnalyticsClient{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *AnalyticsClient) GetMetrics(ctx context.Context, id string) (Metrics, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/users/%s/metrics", c.baseURL, id), nil)
	if err != nil {
		return Metrics{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Metrics{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return Metrics{}, fmt.Errorf("metrics not found")
	}
	if resp.StatusCode != http.StatusOK {
		return Metrics{}, fmt.Errorf("analytics-service returned status %d", resp.StatusCode)
	}

	var response struct {
		Metrics map[string]*float64 `json:"metrics"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return Metrics{}, err
	}

	return Metrics{
		ViewsTotal:        metricInt(response.Metrics, "viewsTotal"),
		FavoritesReceived: metricInt(response.Metrics, "favoritesReceived"),
		DealsClosed:       metricInt(response.Metrics, "dealsClosed"),
		DaysActive:        metricInt(response.Metrics, "daysActive"),
	}, nil
}

func metricInt(metrics map[string]*float64, key string) int {
	value, ok := metrics[key]
	if !ok || value == nil {
		return 0
	}

	return int(*value)
}
