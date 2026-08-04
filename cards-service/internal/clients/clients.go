package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Profile struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Metrics struct {
	ViewsRealty int `json:"views_realty"`
	ViewsAuto   int `json:"views_auto"`
	Purchases   int `json:"purchases"`
	Favorites   int `json:"favorites"`
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

func (c *UserClient) GetProfile(id string) (Profile, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/internal/users/%s", c.baseURL, id))
	if err != nil {
		return Profile{}, err
	}
	defer resp.Body.Close()

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

func (c *AnalyticsClient) GetMetrics(id string) (Metrics, error) {
	resp, err := c.httpClient.Get(fmt.Sprintf("%s/internal/metrics/%s", c.baseURL, id))
	if err != nil {
		return Metrics{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Metrics{}, fmt.Errorf("metrics not found")
	}
	if resp.StatusCode != http.StatusOK {
		return Metrics{}, fmt.Errorf("analytics-service returned status %d", resp.StatusCode)
	}

	var metrics Metrics
	if err := json.NewDecoder(resp.Body).Decode(&metrics); err != nil {
		return Metrics{}, err
	}

	return metrics, nil
}
