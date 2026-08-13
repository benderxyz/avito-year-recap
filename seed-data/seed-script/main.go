package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ingestEvent struct {
	UserID     uint64    `json:"user_id"`
	SessionID  string    `json:"session_id"`
	EventType  string    `json:"event_type"`
	Value      float64   `json:"value"`
	Payload    string    `json:"payload"`
	OccurredAt time.Time `json:"occurred_at"`
}

type userPayload struct {
	ExternalID string `json:"external_id"`
	Username   string `json:"username"`
	Timezone   string `json:"timezone"`
}

type seedUser struct {
	ID       uint64
	Username string
	Activity float64
}

var avitoCategories = []string{
	"electronics",
	"home",
	"transport",
	"fashion",
	"hobby",
	"animals",
	"services",
	"jobs",
	"realty",
}

var defaultUsers = []seedUser{
	{ID: 42, Username: "Alex", Activity: 1.0},
	{ID: 43, Username: "Maria", Activity: 1.6},
	{ID: 44, Username: "Ivan", Activity: 0.45},
	{ID: 45, Username: "Olga", Activity: 1.25},
	{ID: 46, Username: "Dmitry", Activity: 0.75},
	{ID: 47, Username: "Elena", Activity: 1.9},
	{ID: 48, Username: "Sergey", Activity: 0.3},
	{ID: 49, Username: "Anna", Activity: 1.1},
	{ID: 50, Username: "Pavel", Activity: 0.9},
	{ID: 51, Username: "Nina", Activity: 1.4},
}

func main() {
	userFlag := flag.String("user", "", "comma-separated user ids; empty seeds the default roster")
	year := flag.Int("year", 2026, "recap year")
	analyticsURL := flag.String("analytics", "http://localhost:8080", "analytics-service base url")
	userURL := flag.String("users", "http://localhost:8082", "user-service base url")
	seed := flag.Int64("seed", 42, "rng seed")
	flag.Parse()

	users, err := resolveUsers(*userFlag)
	if err != nil {
		fail(err)
	}

	from := time.Date(*year, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(*year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

	totalEvents := 0
	for i, user := range users {
		rng := rand.New(rand.NewSource(*seed + int64(user.ID)*1009 + int64(i)))
		if err := upsertUser(*userURL, user); err != nil {
			fail(err)
		}

		events := generateEvents(user.ID, user.Activity, from, to, rng)
		if err := postEvents(*analyticsURL, events); err != nil {
			fail(err)
		}
		totalEvents += len(events)
		fmt.Printf("seeded user %d (%s) with %d events for %d\n", user.ID, user.Username, len(events), *year)
	}

	fmt.Printf("seeded %d users with %d events total for %d\n", len(users), totalEvents, *year)
}

func resolveUsers(raw string) ([]seedUser, error) {
	if strings.TrimSpace(raw) == "" {
		return defaultUsers, nil
	}

	parts := strings.Split(raw, ",")
	users := make([]seedUser, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.ParseUint(part, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse user id %q: %w", part, err)
		}

		if known, ok := findDefaultUser(id); ok {
			users = append(users, known)
			continue
		}

		users = append(users, seedUser{
			ID:       id,
			Username: fmt.Sprintf("User%d", id),
			Activity: 1.0,
		})
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no user ids provided")
	}
	return users, nil
}

func findDefaultUser(id uint64) (seedUser, bool) {
	for _, user := range defaultUsers {
		if user.ID == id {
			return user, true
		}
	}
	return seedUser{}, false
}

func upsertUser(base string, user seedUser) error {
	body, err := json.Marshal(userPayload{
		ExternalID: fmt.Sprintf("avito-%d", user.ID),
		Username:   user.Username,
		Timezone:   "Europe/Moscow",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%d", base, user.ID), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("upsert user %d: status %d", user.ID, resp.StatusCode)
	}
	return nil
}

func postEvents(base string, events []ingestEvent) error {
	const batchSize = 200
	for start := 0; start < len(events); start += batchSize {
		end := start + batchSize
		if end > len(events) {
			end = len(events)
		}

		body, err := json.Marshal(events[start:end])
		if err != nil {
			return err
		}

		resp, err := http.Post(base+"/events", "application/json", bytes.NewReader(body))
		if err != nil {
			return err
		}
		if resp.StatusCode >= 300 {
			_ = resp.Body.Close()
			return fmt.Errorf("post events batch: status %d", resp.StatusCode)
		}
		_ = resp.Body.Close()
	}
	return nil
}

func generateEvents(userID uint64, activity float64, from, to time.Time, rng *rand.Rand) []ingestEvent {
	if activity <= 0 {
		activity = 0.1
	}

	events := make([]ingestEvent, 0, 2048)
	session := uuid.NewString()

	firstPublish := from.AddDate(0, 0, 8+rng.Intn(20)).Add(time.Duration(8+rng.Intn(8)) * time.Hour)
	events = append(events,
		evt(userID, session, "first_item_published", 0, "{}", firstPublish),
		evt(userID, session, "item_published", 1, "{}", firstPublish.Add(time.Minute)),
		evt(userID, session, "active_items_count", 1, "{}", firstPublish.Add(2*time.Minute)),
		evt(userID, session, "day_active", 0, "{}", firstPublish),
	)

	activeItems := 1.0
	peakDayViews := 0.0
	firstDealDone := false
	dayActiveChance := clampFloat(0.2+0.35*activity, 0.15, 0.85)
	publishChance := clampFloat(0.12+0.16*activity, 0.05, 0.5)
	dealChance := clampFloat(0.08+0.12*activity, 0.03, 0.4)

	day := from
	for day.Before(to) {
		if rng.Float64() > dayActiveChance {
			day = day.AddDate(0, 0, 1)
			continue
		}

		session = uuid.NewString()
		at := day.Add(time.Duration(9+rng.Intn(10)) * time.Hour).Add(time.Duration(rng.Intn(60)) * time.Minute)
		events = append(events, evt(userID, session, "day_active", 0, "{}", at))

		searches := 1 + rng.Intn(1+int(3*activity))
		for i := 0; i < searches; i++ {
			category := avitoCategories[rng.Intn(len(avitoCategories))]
			events = append(events,
				evt(userID, session, "search_performed", 1, "{}", at.Add(time.Duration(i)*3*time.Minute)),
				evt(userID, session, "category_opened", 0, payloadCategory(category), at.Add(time.Duration(i)*3*time.Minute+time.Minute)),
			)
		}

		if rng.Float64() < publishChance {
			activeItems++
			events = append(events,
				evt(userID, session, "item_published", 1, "{}", at.Add(20*time.Minute)),
				evt(userID, session, "active_items_count", activeItems, "{}", at.Add(21*time.Minute)),
			)
		}

		dayViews := float64(10+rng.Intn(40)) + float64(rng.Intn(1+int(160*activity)))
		events = append(events, evt(userID, session, "item_view_received", dayViews, "{}", at.Add(30*time.Minute)))
		if dayViews > peakDayViews {
			peakDayViews = dayViews
			events = append(events, evt(userID, session, "peak_day_views", peakDayViews, "{}", at.Add(31*time.Minute)))
		}

		favs := float64(rng.Intn(1 + int(8*activity)))
		if favs > 0 {
			events = append(events, evt(userID, session, "item_favorited", favs, "{}", at.Add(40*time.Minute)))
		}

		messages := 1 + rng.Intn(1+int(6*activity))
		for i := 0; i < messages; i++ {
			msgAt := at.Add(time.Duration(50+i*7) * time.Minute)
			events = append(events, evt(userID, session, "chat_message_sent", 1, "{}", msgAt))
			if rng.Float64() < 0.7 {
				startUnix := float64(msgAt.Unix())
				replyAt := msgAt.Add(time.Duration(2+rng.Intn(40)) * time.Minute)
				events = append(events, evt(userID, session, "chat_reply", startUnix, "{}", replyAt))
			}
		}

		if rng.Float64() < dealChance && activeItems >= 1 {
			dealAt := at.Add(3 * time.Hour)
			sale := float64(1500 + rng.Intn(45000))
			events = append(events,
				evt(userID, session, "deal_completed", 1, "{}", dealAt),
				evt(userID, session, "sale_proceeds", sale, "{}", dealAt.Add(time.Minute)),
				evt(userID, session, "buyer_saving", float64(200+rng.Intn(5000)), "{}", dealAt.Add(2*time.Minute)),
			)
			if !firstDealDone {
				events = append(events, evt(userID, session, "first_deal_completed", 0, "{}", dealAt))
				firstDealDone = true
			}
			if rng.Float64() < 0.45 {
				events = append(events, evt(userID, session, "delivery_order_created", 1, "{}", dealAt.Add(5*time.Minute)))
			}
			activeItems = maxFloat(0, activeItems-1)
			events = append(events, evt(userID, session, "active_items_count", activeItems, "{}", dealAt.Add(6*time.Minute)))
		}

		day = day.AddDate(0, 0, 1)
	}

	return events
}

func evt(userID uint64, sessionID, eventType string, value float64, payload string, at time.Time) ingestEvent {
	if payload == "" {
		payload = "{}"
	}
	return ingestEvent{
		UserID:     userID,
		SessionID:  sessionID,
		EventType:  eventType,
		Value:      value,
		Payload:    payload,
		OccurredAt: at.UTC(),
	}
}

func payloadCategory(category string) string {
	raw, _ := json.Marshal(map[string]string{"category": category})
	return string(raw)
}

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "seed-script: %v\n", err)
	os.Exit(1)
}
