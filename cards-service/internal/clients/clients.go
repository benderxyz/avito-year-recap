package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type Profile struct {
	UserID     uint64 `json:"user_id"`
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}

type MetricsResponse struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	UserID   uint64  `json:"user_id"`
	Timezone string  `json:"timezone"`
	Metrics  Metrics `json:"metrics"`
}

type Metrics struct {
	ActiveListings    int     `json:"activeListings"`
	AvgReplySeconds   float64 `json:"avgReplySeconds"`
	CategoriesTried   int     `json:"categoriesTried"`
	DaysActive        int     `json:"daysActive"`
	DealsClosed       int     `json:"dealsClosed"`
	DeliveryOrders    int     `json:"deliveryOrders"`
	FavoritesReceived int     `json:"favoritesReceived"`
	FirstDealAt       int64   `json:"firstDealAt"`
	FirstListingAt    int64   `json:"firstListingAt"`
	ListingsPublished int     `json:"listingsPublished"`
	MessagesSent      int     `json:"messagesSent"`
	MoneyEarned       int64   `json:"moneyEarned"`
	MoneySaved        int64   `json:"moneySaved"`
	PeakDayViews      int     `json:"peakDayViews"`
	SearchQueries     int     `json:"searchQueries"`
	SellerRating      float64 `json:"sellerRating"`
	ViewsTotal        int     `json:"viewsTotal"`
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
	url := fmt.Sprintf("%s/users/%s", c.baseURL, id)

	slog.Debug("requesting user profile",
		"user_id", id,
		"url", url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		slog.Error("failed to create user profile request",
			"error", err,
		)
		return Profile{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("user profile request failed",
			"error", err,
		)
		return Profile{}, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	slog.Debug("user profile response received",
		"status", resp.StatusCode,
	)

	if resp.StatusCode == http.StatusNotFound {
		return Profile{}, fmt.Errorf("profile not found")
	}

	if resp.StatusCode != http.StatusOK {
		return Profile{}, fmt.Errorf(
			"user-service returned status %d",
			resp.StatusCode,
		)
	}

	var profile Profile

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		slog.Error("failed to decode user profile response",
			"error", err,
		)
		return Profile{}, err
	}

	slog.Info("user profile loaded",
		"user_id", id,
	)

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

func (c *AnalyticsClient) Health(ctx context.Context) error {
	url := c.baseURL + "/health"

	slog.Debug("checking analytics service health",
		"url", url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"analytics service unhealthy: %s",
			resp.Status,
		)
	}

	return nil
}

func (c *AnalyticsClient) GetMetrics(
	ctx context.Context,
	id string,
	year int,
) (Metrics, error) {
	if err := c.Health(ctx); err != nil {
		slog.Error("analytics health check failed",
			"error", err,
		)
		return Metrics{}, err
	}

	slog.Debug("analytics service is healthy")

	url := fmt.Sprintf(
		"%s/users/%s/metrics",
		c.baseURL,
		id,
	)

	slog.Debug("requesting user metrics",
		"user_id", id,
		"year", year,
		"url", url,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		slog.Error("failed to create metrics request",
			"error", err,
		)
		return Metrics{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("metrics request failed",
			"error", err,
		)
		return Metrics{}, err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}()

	slog.Debug("metrics response received",
		"status", resp.StatusCode,
	)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		slog.Error("analytics service returned non-success response",
			"status", resp.StatusCode,
			"body", string(body),
		)

		return Metrics{}, fmt.Errorf(
			"analytics-service returned status %d",
			resp.StatusCode,
		)
	}

	var response MetricsResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.Error("failed to decode metrics response",
			"error", err,
		)
		return Metrics{}, err
	}

	slog.Info("metrics loaded successfully",
		"user_id", id,
	)

	return response.Metrics, nil
}
