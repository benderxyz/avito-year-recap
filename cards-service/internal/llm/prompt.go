package llm

import (
	"fmt"
	"math"
	"strconv"
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

func buildMessages(displayName string, badges []models.Badge, m clients.Metrics, defs []models.MetricDefinition) []Message {
	var b strings.Builder

	name := sanitizeName(displayName)
	fmt.Fprintf(&b, "Пользователь: %s\n\n", name)

	b.WriteString("Отобранные бейджи (перепиши каждый, сохранив id):\n")
	for _, badge := range badges {
		fmt.Fprintf(&b, "- id=%s | title=%q | description=%q\n", badge.ID, badge.Title, badge.Description)
	}

	b.WriteString("\nМетрики пользователя за год (для рассуждения об инсайте, не выводи их числами):\n")
	for _, line := range metricsSummary(m, defs) {
		fmt.Fprintf(&b, "- %s\n", line)
	}

	return []Message{
		{Role: RoleSystem, Content: systemPrompt},
		{Role: RoleUser, Content: b.String()},
	}
}

func sanitizeName(displayName string) string {
	name := whitespaceRe.ReplaceAllString(strings.TrimSpace(displayName), " ")
	if name == "" {
		return "пользователь"
	}
	if runeLen(name) > maxDisplayNameLen {
		name = string([]rune(name)[:maxDisplayNameLen])
	}
	return name
}

func metricsSummary(m clients.Metrics, defs []models.MetricDefinition) []string {
	lines := make([]string, 0, len(defs))
	for _, def := range defs {
		if !def.IncludeInLLM || !def.IsPublic {
			continue
		}

		value, ok := m.Value(def)
		if !ok || value == 0 {
			continue
		}

		lines = append(lines, fmt.Sprintf("%s: %s", def.Key, formatPromptNumber(value)))
	}
	return lines
}

func formatPromptNumber(value float64) string {
	if value == math.Trunc(value) {
		return strconv.FormatInt(int64(value), 10)
	}
	return strconv.FormatFloat(value, 'f', 1, 64)
}
