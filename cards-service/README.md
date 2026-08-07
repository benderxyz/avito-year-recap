# Cards Service

Overview

The Cards Service prepares "recap" cards for a user: it fetches the user profile from the user-service and metrics from the analytics-service, then generates a set of cards returned to the client.

Key files

- [entrypoint: main.go](/home/darina/study/avito_hack/avito-year-recap.worktrees/seed-data-users-api-user-service/cards-service/cmd/main.go)
- [HTTP handlers](/home/darina/study/avito_hack/avito-year-recap.worktrees/seed-data-users-api-user-service/cards-service/internal/api/handlers.go)
- [Service clients (user / analytics)](/home/darina/study/avito_hack/avito-year-recap.worktrees/seed-data-users-api-user-service/cards-service/internal/clients/clients.go)
- [Card generator](/home/darina/study/avito_hack/avito-year-recap.worktrees/seed-data-users-api-user-service/cards-service/internal/cards/generator.go)
- [Recap models](/home/darina/study/avito_hack/avito-year-recap.worktrees/seed-data-users-api-user-service/cards-service/internal/models/recap.go)

API

- GET /health
  - Description: service health check
  - Response: text/plain, body: "cards-service: OK\n"

- GET /api/recap/{id}
  - Description: build a recap for the profile with the provided id
  - Behavior: requests profile from user-service and metrics from analytics-service, then invokes the card generator
  - Success response: 200 OK, application/json
  - Response shape (Recap):

    {
      "profile_id": "<id>",
      "cards": [
        {
          "type": "<type>",
          "title": "<title>",
          "text": "<text>",
          "action": "<optional>",
          "action_value": "<optional>"
        }
      ]
    }

  - Errors:
    - 404 "profile not found" — when user-service returns NotFound
    - 404 "metrics not found" — when analytics-service returns an error

Configuration

The service is configured via environment variables (see [cmd/main.go](/home/darina/study/avito_hack/avito-year-recap.worktrees/seed-data-users-api-user-service/cards-service/cmd/main.go)):

- USER_SERVICE_URL — URL of the user-service (default: http://localhost:8082)
- ANALYTICS_SERVICE_URL — URL of the analytics-service (default: http://localhost:8080)
- The HTTP server address is configured in code (default: :8081)


