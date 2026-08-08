package clients

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"analytics-service/internal/models"
)

var ErrUserNotFound = errors.New("user not found")

type UserClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserClient(baseURL string) *UserClient {
	return &UserClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *UserClient) GetUser(ctx context.Context, userID uint64) (models.User, error) {
	requestURL := fmt.Sprintf("%s/users/%d", c.baseURL, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return models.User{}, fmt.Errorf("create user request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return models.User{}, fmt.Errorf("user-service request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return models.User{}, ErrUserNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return models.User{}, fmt.Errorf("user-service returned status %d", resp.StatusCode)
	}

	var user models.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return models.User{}, fmt.Errorf("decode user response: %w", err)
	}

	return user, nil
}

func (c *UserClient) Timezone(ctx context.Context, userID uint64) (string, error) {
	user, err := c.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return "UTC", nil
		}
		return "", fmt.Errorf("resolve timezone: %w", err)
	}

	timezone := strings.TrimSpace(user.Timezone)
	if timezone == "" {
		return "UTC", nil
	}

	return timezone, nil
}
