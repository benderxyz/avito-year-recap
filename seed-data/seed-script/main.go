package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
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

func main() {
	userID := flag.Uint64("user", 42, "user id")
	year := flag.Int("year", 2026, "recap year")
	analyticsURL := flag.String("analytics", "http://127.0.0.1:8080", "analytics-service base url")
	//userURL := flag.String("users", "http://127.0.0.1:8082", "user-service base url")
	seed := flag.Int64("seed", 42, "rng seed")
	flag.Parse()

	rng := rand.New(rand.NewSource(*seed))
	from := time.Date(*year, time.January, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(*year+1, time.January, 1, 0, 0, 0, 0, time.UTC)

	// user-service пока не интегрирован в seed-script, поэтому закомментировал вызов upsertUser
	/*if err := upsertUser(*userURL, *userID); err != nil {
		fail(err)
	}*/

	events := generateEvents(*userID, from, to, rng)
	if err := postEvents(*analyticsURL, events); err != nil {
		fail(err)
	}

	fmt.Printf("seeded user %d with %d events for %d\n", *userID, len(events), *year)
}

func upsertUser(base string, userID uint64) error {
	body, err := json.Marshal(userPayload{
		ExternalID: fmt.Sprintf("avito-%d", userID),
		Username:   "Alex",
		Timezone:   "Europe/Moscow",
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%d", base, userID), bytes.NewReader(body))
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
		return fmt.Errorf("upsert user: status %d", resp.StatusCode)
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

func generateEvents(userID uint64, from, to time.Time, rng *rand.Rand) []ingestEvent {
	events := make([]ingestEvent, 0, 2048)
	session := uuid.NewString()

	firstPublish := from.AddDate(0, 0, 12).Add(10 * time.Hour)
	events = append(events,
		evt(userID, session, "first_item_published", 0, "{}", firstPublish),
		evt(userID, session, "item_published", 1, "{}", firstPublish.Add(time.Minute)),
		evt(userID, session, "active_items_count", 1, "{}", firstPublish.Add(2*time.Minute)),
		evt(userID, session, "day_active", 0, "{}", firstPublish),
	)

	activeItems := 1.0
	peakDayViews := 0.0
	sellerRating := 4.2
	firstDealDone := false

	day := from
	for day.Before(to) {
		if rng.Float64() > 0.62 {
			day = day.AddDate(0, 0, 1)
			continue
		}

		session = uuid.NewString()
		at := day.Add(time.Duration(9+rng.Intn(10)) * time.Hour).Add(time.Duration(rng.Intn(60)) * time.Minute)
		events = append(events, evt(userID, session, "day_active", 0, "{}", at))

		searches := 1 + rng.Intn(4)
		for i := 0; i < searches; i++ {
			category := avitoCategories[rng.Intn(len(avitoCategories))]
			events = append(events,
				evt(userID, session, "search_performed", 1, "{}", at.Add(time.Duration(i)*3*time.Minute)),
				evt(userID, session, "category_opened", 0, payloadCategory(category), at.Add(time.Duration(i)*3*time.Minute+time.Minute)),
			)
		}

		if rng.Float64() < 0.28 {
			activeItems++
			events = append(events,
				evt(userID, session, "item_published", 1, "{}", at.Add(20*time.Minute)),
				evt(userID, session, "active_items_count", activeItems, "{}", at.Add(21*time.Minute)),
			)
		}

		dayViews := float64(20 + rng.Intn(180))
		events = append(events, evt(userID, session, "item_view_received", dayViews, "{}", at.Add(30*time.Minute)))
		if dayViews > peakDayViews {
			peakDayViews = dayViews
			events = append(events, evt(userID, session, "peak_day_views", peakDayViews, "{}", at.Add(31*time.Minute)))
		}

		favs := float64(rng.Intn(8))
		if favs > 0 {
			events = append(events, evt(userID, session, "item_favorited", favs, "{}", at.Add(40*time.Minute)))
		}

		messages := 1 + rng.Intn(6)
		for i := 0; i < messages; i++ {
			msgAt := at.Add(time.Duration(50+i*7) * time.Minute)
			events = append(events, evt(userID, session, "chat_message_sent", 1, "{}", msgAt))
			if rng.Float64() < 0.7 {
				startUnix := float64(msgAt.Unix())
				replyAt := msgAt.Add(time.Duration(2+rng.Intn(40)) * time.Minute)
				events = append(events, evt(userID, session, "chat_reply", startUnix, "{}", replyAt))
			}
		}

		if rng.Float64() < 0.18 && activeItems >= 1 {
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
			sellerRating = minFloat(5, sellerRating+0.02)
			events = append(events, evt(userID, session, "seller_rating", sellerRating, "{}", dealAt.Add(7*time.Minute)))
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

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "seed-script: %v\n", err)
	os.Exit(1)
}
