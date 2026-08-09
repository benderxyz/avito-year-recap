package llm

import (
	"fmt"
	"strings"

	"cards-service/internal/clients"
	"cards-service/internal/models"
)

const systemPrompt = `Ты — копирайтер сервиса «Итоги года на Авито». Пишешь на русском языке (ru-RU).

Твоя задача:
1. Переписать тексты уже отобранных бейджей (достижений) — живее и теплее, сохранив смысл.
2. Придумать один короткий позитивный инсайт о пользователе на основе комбинации его метрик
   (например: мало переписок, но много покупок → вдумчивый охотник за находками).

Строгие правила:
- Только позитивный, дружелюбный тон. Никакой агрессии, критики, сарказма, обесценивания.
- НЕ используй конкретику: никаких цифр, сумм, валют, процентов, названий городов.
  Числа и суммы фронтенд показывает отдельно — твой текст должен читаться без них.
- Пиши кратко. Заголовок бейджа — до 40 символов, описание — до 120. Инсайт — заголовок до 60,
  текст до 200 символов. Каждое поле — одна строка, без переносов.
- Не выдумывай факты, которых нет в метриках.

Ответь СТРОГО одним JSON-объектом, без markdown и пояснений, по схеме:
{
  "badges": { "<id бейджа>": { "title": "...", "description": "..." } },
  "insight": { "title": "...", "text": "..." }
}
Ключи в "badges" — это id из списка бейджей ниже. Если инсайта нет — верни "insight": null.`

func buildMessages(displayName string, badges []models.Badge, m clients.Metrics) []Message {
	var b strings.Builder

	name := displayName
	if name == "" {
		name = "пользователь"
	}
	fmt.Fprintf(&b, "Пользователь: %s\n\n", name)

	b.WriteString("Отобранные бейджи (перепиши каждый, сохранив id):\n")
	for _, badge := range badges {
		fmt.Fprintf(&b, "- id=%s | title=%q | description=%q\n", badge.ID, badge.Title, badge.Description)
	}

	b.WriteString("\nМетрики пользователя за год (для рассуждения об инсайте, не выводи их числами):\n")
	for _, line := range metricsSummary(m) {
		fmt.Fprintf(&b, "- %s\n", line)
	}

	return []Message{
		{Role: RoleSystem, Content: systemPrompt},
		{Role: RoleUser, Content: b.String()},
	}
}

func metricsSummary(m clients.Metrics) []string {
	pairs := []struct {
		key string
		val int64
	}{
		{"listingsPublished", int64(m.ListingsPublished)},
		{"viewsTotal", int64(m.ViewsTotal)},
		{"favoritesReceived", int64(m.FavoritesReceived)},
		{"messagesSent", int64(m.MessagesSent)},
		{"dealsClosed", int64(m.DealsClosed)},
		{"moneyEarned", m.MoneyEarned},
		{"moneySaved", m.MoneySaved},
		{"daysActive", int64(m.DaysActive)},
		{"peakDayViews", int64(m.PeakDayViews)},
		{"searchQueries", int64(m.SearchQueries)},
		{"categoriesTried", int64(m.CategoriesTried)},
		{"deliveryOrders", int64(m.DeliveryOrders)},
		{"activeListings", int64(m.ActiveListings)},
	}

	lines := make([]string, 0, len(pairs))
	for _, p := range pairs {
		if p.val != 0 {
			lines = append(lines, fmt.Sprintf("%s: %d", p.key, p.val))
		}
	}
	if m.SellerRating != 0 {
		lines = append(lines, fmt.Sprintf("sellerRating: %.1f", m.SellerRating))
	}
	return lines
}
