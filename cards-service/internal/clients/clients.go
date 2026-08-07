package clients

import (
	"cards-service/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

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
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

}

func (c *UserClient) GetProfile(ctx context.Context, id string) (models.Profile, error) {
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
		return models.Profile{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("user profile request failed",
			"error", err,
		)
		return models.Profile{}, err
	}

	defer resp.Body.Close()

	slog.Debug("user profile response received",
		"status", resp.StatusCode,
	)

	if resp.StatusCode == http.StatusNotFound {
		return models.Profile{}, fmt.Errorf("profile not found")
	}

	if resp.StatusCode != http.StatusOK {
		return models.Profile{}, fmt.Errorf(
			"user-service returned status %d",
			resp.StatusCode,
		)
	}

	var profile models.Profile

	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		slog.Error("failed to decode user profile response",
			"error", err,
		)
		return models.Profile{}, err
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

	defer resp.Body.Close()
	return nil
}

func (c *AnalyticsClient) GetMetrics(
	ctx context.Context,
	id string,
	year int,
	timezone string,
) (Metrics, error) {

	from, to, err := getYearRange(year, timezone)
	if err != nil {
		return Metrics{}, err
	}

	query := url.Values{}
	query.Set("from", from)
	query.Set("to", to)

	requestURL := fmt.Sprintf(
		"%s/users/%s/metrics?%s",
		c.baseURL,
		id,
		query.Encode(),
	)

	slog.Debug("requesting user metrics",
		"user_id", id,
		"year", year,
		"url", requestURL,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		requestURL,
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

	defer resp.Body.Close()

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

func getYearRange(year int, timezone string) (string, string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", "", fmt.Errorf(
			"invalid timezone %s: %w",
			timezone,
			err,
		)
	}

	from := time.Date(
		year,
		time.January,
		1,
		0,
		0,
		0,
		0,
		loc,
	)

	to := time.Date(
		year,
		time.December,
		31,
		23,
		59,
		59,
		0,
		loc,
	)

	return from.Format(time.RFC3339), to.Format(time.RFC3339), nil

}
