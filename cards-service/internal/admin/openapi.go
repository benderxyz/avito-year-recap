package admin

import (
	"net/http"
	"reflect"
	"strings"
)

func (h *Handler) OpenAPI(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, openAPIDocument())
}

func (h *Handler) Docs(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if _, err := w.Write([]byte(swaggerUIPage)); err != nil {
		http.Error(w, "failed to write docs page", http.StatusInternalServerError)
	}
}

const swaggerUIPage = `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>cards-service admin API</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: '/api/admin/openapi.json',
      dom_id: '#swagger-ui',
      persistAuthorization: true
    });
  </script>
</body>
</html>
`

var propertyEnums = map[string][]string{
	"valueType":   asStrings(valueTypes),
	"sourceField": asStrings(sourceFields),
	"visibility":  asStrings(visibilities),
	"sceneType":   asStrings(sceneTypes),
	"op":          asStrings(predicateOps),
	"match":       asStrings(matchModes),
	"currency":    {"RUB"},
}

var requiredFields = map[string][]string{
	"MetricDefinition":     {"key", "valueType", "isPublic", "sourceKey", "sourceField", "includeInLlm", "sortOrder", "enabled", "createdAt", "updatedAt"},
	"MetricWrite":          {"valueType"},
	"MetricCreate":         {"key", "valueType"},
	"BadgeRule":            {"id", "title", "description", "visibility", "when", "sortOrder", "enabled", "createdAt", "updatedAt"},
	"BadgeWrite":           {"title", "description", "visibility", "when"},
	"BadgeCreate":          {"id", "title", "description", "visibility", "when"},
	"StoryRule":            {"id", "sceneType", "visibility", "payload", "sortOrder", "enabled", "createdAt", "updatedAt"},
	"StoryWrite":           {"sceneType", "visibility", "payload"},
	"StoryCreate":          {"id", "sceneType", "visibility", "payload"},
	"RecommendationRule":   {"id", "title", "text", "callout", "linkLabel", "path", "priority", "when", "enabled", "createdAt", "updatedAt"},
	"RecommendationWrite":  {"title", "text", "linkLabel", "path", "when"},
	"RecommendationCreate": {"id", "title", "text", "linkLabel", "path", "when"},
	"Predicate":            {"metric", "op"},
	"GroupWhen":            {"predicates"},
	"AdminError":           {"error"},
}

func asStrings[T ~string](values []T) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		result = append(result, string(value))
	}
	return result
}

type schemaBuilder struct {
	schemas map[string]any
}

func (b *schemaBuilder) ref(t reflect.Type) map[string]any {
	name := t.Name()
	if _, known := b.schemas[name]; !known {
		b.schemas[name] = map[string]any{"type": "object"}
		b.schemas[name] = b.object(t)
	}

	return map[string]any{"$ref": "#/components/schemas/" + name}
}

func (b *schemaBuilder) object(t reflect.Type) map[string]any {
	properties := map[string]any{}
	b.collect(t, properties)

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}
	if required, ok := requiredFields[t.Name()]; ok {
		schema["required"] = required
	}

	return schema
}

func (b *schemaBuilder) collect(t reflect.Type, properties map[string]any) {
	for index := 0; index < t.NumField(); index++ {
		field := t.Field(index)

		if field.Anonymous && field.Type.Kind() == reflect.Struct && field.Tag.Get("json") == "" {
			b.collect(field.Type, properties)
			continue
		}

		name := strings.Split(field.Tag.Get("json"), ",")[0]
		if name == "" || name == "-" {
			continue
		}

		properties[name] = b.field(field.Type, name)
	}
}

func (b *schemaBuilder) field(t reflect.Type, name string) map[string]any {
	if t.Kind() == reflect.Pointer {
		schema := b.field(t.Elem(), name)
		if ref, isRef := schema["$ref"]; isRef {
			return map[string]any{
				"nullable": true,
				"allOf":    []any{map[string]any{"$ref": ref}},
			}
		}
		schema["nullable"] = true
		return schema
	}

	switch t.Kind() {
	case reflect.String:
		schema := map[string]any{"type": "string"}
		if values, ok := propertyEnums[name]; ok {
			schema["enum"] = values
		}
		return schema
	case reflect.Bool:
		return map[string]any{"type": "boolean"}
	case reflect.Int, reflect.Int64:
		return map[string]any{"type": "integer"}
	case reflect.Float64:
		return map[string]any{"type": "number"}
	case reflect.Slice:
		return map[string]any{"type": "array", "items": b.field(t.Elem(), name)}
	case reflect.Map:
		return map[string]any{"type": "object", "additionalProperties": true}
	case reflect.Struct:
		return b.ref(t)
	case reflect.Interface:
		return map[string]any{}
	default:
		return map[string]any{}
	}
}

type queryParam struct {
	name        string
	description string
	schema      map[string]any
}

func enumParam(name, description string, values []string) queryParam {
	return queryParam{
		name:        name,
		description: description,
		schema:      map[string]any{"type": "string", "enum": values},
	}
}

func boolParam(name, description string) queryParam {
	return queryParam{name: name, description: description, schema: map[string]any{"type": "boolean"}}
}

func stringParam(name, description string) queryParam {
	return queryParam{name: name, description: description, schema: map[string]any{"type": "string"}}
}

func intParam(name, description string) queryParam {
	return queryParam{name: name, description: description, schema: map[string]any{"type": "integer"}}
}

type collection struct {
	path      string
	tag       string
	param     string
	item      any
	create    any
	write     any
	filters   []queryParam
	conflicts bool
}

func openAPIDocument() map[string]any {
	builder := &schemaBuilder{schemas: map[string]any{}}

	collections := []collection{
		{
			path:   "metrics",
			tag:    "metrics",
			param:  "key",
			item:   MetricDefinition{},
			create: MetricCreate{},
			write:  MetricWrite{},
			filters: []queryParam{
				boolParam("enabled", "Только включенные или только выключенные метрики"),
				boolParam("isPublic", "Метрики, доступные в публичном recap"),
				boolParam("includeInLlm", "Метрики, которые уходят в промпт LLM"),
				enumParam("valueType", "Тип значения метрики", asStrings(valueTypes)),
				enumParam("sourceField", "Поле ответа analytics-service", asStrings(sourceFields)),
				stringParam("search", "Подстрока ключа метрики"),
			},
			conflicts: true,
		},
		{
			path:   "badges",
			tag:    "badges",
			param:  "id",
			item:   BadgeRule{},
			create: BadgeCreate{},
			write:  BadgeWrite{},
			filters: []queryParam{
				boolParam("enabled", "Только включенные или только выключенные бейджи"),
				enumParam("visibility", "Режим recap, в котором виден бейдж", asStrings(visibilities)),
				stringParam("metric", "Ключ метрики в условии"),
				stringParam("search", "Подстрока id, заголовка или описания"),
			},
			conflicts: true,
		},
		{
			path:   "stories",
			tag:    "stories",
			param:  "id",
			item:   StoryRule{},
			create: StoryCreate{},
			write:  StoryWrite{},
			filters: []queryParam{
				boolParam("enabled", "Только включенные или только выключенные сцены"),
				enumParam("visibility", "Режим recap, в котором видна сцена", asStrings(visibilities)),
				enumParam("sceneType", "Тип сцены", asStrings(sceneTypes)),
				stringParam("metric", "Ключ метрики в условии"),
				stringParam("search", "Подстрока id или payload"),
			},
		},
		{
			path:   "recommendations",
			tag:    "recommendations",
			param:  "id",
			item:   RecommendationRule{},
			create: RecommendationCreate{},
			write:  RecommendationWrite{},
			filters: []queryParam{
				boolParam("enabled", "Только включенные или только выключенные рекомендации"),
				stringParam("metric", "Ключ метрики в предикатах"),
				intParam("minPriority", "Минимальный приоритет"),
				stringParam("search", "Подстрока id, заголовка или текста"),
			},
		},
	}

	paths := map[string]any{}
	for _, entry := range collections {
		addCollectionPaths(paths, builder, entry)
	}

	paths["/api/admin/preview"] = map[string]any{
		"get": map[string]any{
			"tags":        []string{"preview"},
			"summary":     "Recap на одноразовых тестовых данных, без LLM",
			"description": "Метрики генерируются случайно по текущим определениям, поэтому ответ каждый раз разный. Параметр seed делает результат воспроизводимым.",
			"parameters": append(
				pathParameters(nil),
				queryParameters([]queryParam{
					intParam("year", "Год recap, по умолчанию текущий"),
					enumParam("mode", "Режим recap", []string{"private", "public"}),
					intParam("seed", "Зерно генератора случайных данных"),
				})...,
			),
			"responses": map[string]any{
				"200": map[string]any{
					"description": "RecapPayload",
					"content": map[string]any{
						"application/json": map[string]any{
							"schema": map[string]any{"type": "object"},
						},
					},
				},
				"400": errorResponseSpec("Некорректный параметр запроса"),
				"401": errorResponseSpec("Нет токена или токен неверный"),
				"404": errorResponseSpec("В админке нет ни одной метрики"),
			},
		},
	}

	builder.ref(reflect.TypeOf(AdminError{}))

	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "cards-service admin API",
			"version":     "1.0.0",
			"description": "Каталоги правил recap: метрики, бейджи, сцены и рекомендации. Все запросы, кроме этого описания, требуют заголовок Authorization: Bearer <ADMIN_API_TOKEN>.",
		},
		"servers": []any{map[string]any{"url": "/"}},
		"security": []any{
			map[string]any{"bearerAuth": []any{}},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":   "http",
					"scheme": "bearer",
				},
			},
			"schemas": builder.schemas,
		},
		"paths": paths,
	}
}

func addCollectionPaths(paths map[string]any, builder *schemaBuilder, entry collection) {
	itemRef := builder.ref(reflect.TypeOf(entry.item))
	createRef := builder.ref(reflect.TypeOf(entry.create))
	writeRef := builder.ref(reflect.TypeOf(entry.write))

	deleteResponses := map[string]any{
		"204": map[string]any{"description": "Запись удалена"},
		"401": errorResponseSpec("Нет токена или токен неверный"),
		"404": errorResponseSpec("Записи нет"),
	}
	if entry.conflicts {
		deleteResponses["409"] = errorResponseSpec("На запись ссылаются другие правила")
	}

	paths["/api/admin/"+entry.path] = map[string]any{
		"get": map[string]any{
			"tags":       []string{entry.tag},
			"summary":    "Список без пагинации, включая записи с enabled = false",
			"parameters": queryParameters(entry.filters),
			"responses": map[string]any{
				"200": jsonResponse("Список записей", map[string]any{
					"type": "object",
					"properties": map[string]any{
						"items": map[string]any{"type": "array", "items": itemRef},
					},
					"required": []string{"items"},
				}),
				"400": errorResponseSpec("Некорректный параметр фильтра"),
				"401": errorResponseSpec("Нет токена или токен неверный"),
			},
		},
		"post": map[string]any{
			"tags":        []string{entry.tag},
			"summary":     "Создать запись",
			"requestBody": jsonBody(createRef),
			"responses": map[string]any{
				"201": jsonResponse("Запись создана", itemRef),
				"400": errorResponseSpec("Тело не прошло валидацию"),
				"401": errorResponseSpec("Нет токена или токен неверный"),
				"409": errorResponseSpec("Запись с таким идентификатором уже есть"),
			},
		},
	}

	paths["/api/admin/"+entry.path+"/{"+entry.param+"}"] = map[string]any{
		"get": map[string]any{
			"tags":       []string{entry.tag},
			"summary":    "Одна запись",
			"parameters": pathParameters([]string{entry.param}),
			"responses": map[string]any{
				"200": jsonResponse("Запись", itemRef),
				"401": errorResponseSpec("Нет токена или токен неверный"),
				"404": errorResponseSpec("Записи нет"),
			},
		},
		"put": map[string]any{
			"tags":        []string{entry.tag},
			"summary":     "Обновить существующую запись",
			"description": "Идентификатор в теле необязателен, но если он есть, он должен совпадать с идентификатором из пути.",
			"parameters":  pathParameters([]string{entry.param}),
			"requestBody": jsonBody(writeRef),
			"responses": map[string]any{
				"200": jsonResponse("Запись обновлена", itemRef),
				"400": errorResponseSpec("Тело не прошло валидацию"),
				"401": errorResponseSpec("Нет токена или токен неверный"),
				"404": errorResponseSpec("Записи нет"),
			},
		},
		"delete": map[string]any{
			"tags":       []string{entry.tag},
			"summary":    "Удалить запись",
			"parameters": pathParameters([]string{entry.param}),
			"responses":  deleteResponses,
		},
	}
}

func pathParameters(names []string) []any {
	parameters := make([]any, 0, len(names))
	for _, name := range names {
		parameters = append(parameters, map[string]any{
			"name":     name,
			"in":       "path",
			"required": true,
			"schema":   map[string]any{"type": "string"},
		})
	}
	return parameters
}

func queryParameters(params []queryParam) []any {
	parameters := make([]any, 0, len(params))
	for _, param := range params {
		parameters = append(parameters, map[string]any{
			"name":        param.name,
			"in":          "query",
			"required":    false,
			"description": param.description,
			"schema":      param.schema,
		})
	}
	return parameters
}

func jsonBody(schema map[string]any) map[string]any {
	return map[string]any{
		"required": true,
		"content": map[string]any{
			"application/json": map[string]any{"schema": schema},
		},
	}
}

func jsonResponse(description string, schema map[string]any) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{"schema": schema},
		},
	}
}

func errorResponseSpec(description string) map[string]any {
	return jsonResponse(description, map[string]any{"$ref": "#/components/schemas/AdminError"})
}
