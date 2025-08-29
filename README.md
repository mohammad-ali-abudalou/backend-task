# Backend Task (Go + Gin + GORM + Postgres)

Implements the user management API with automatic, capacity-aware age-group assignment.

## Features
- REST API (Gin)
- ORM (GORM) with Postgres (tests use SQLite)
- Concurrency-safe group allocation via dedicated `groups` table + row locks
- Validation (email format, uniqueness, DOB in the past)
- Unit tests
- Docker & docker-compose

## Run (Docker Compose)
```bash
docker compose up --build
# API at http://localhost:8080
```

## Run locally (no Docker)
```bash
export DB_DRIVER=postgres
export DB_DSN="host=localhost user=postgres password=postgres dbname=backend port=5432 sslmode=disable TimeZone=UTC"
# start Postgres first, then:
go run ./cmd
```

> **Note for Windows (PowerShell):** run `go run ./cmd` (not `go run cmd/main.go`) to ensure module imports resolve correctly.

## Endpoints

### Create user
`POST /users`
```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "date_of_birth": "1990-05-10"
}
```
**201 Created**
```json
{
  "id": "<uuid>",
  "name": "Alice",
  "email": "alice@example.com",
  "date_of_birth": "1990-05-10T00:00:00Z",
  "group": "adult-1",
  "created_at": "...",
  "updated_at": "..."
}
```

### Get user by id
`GET /users/{id}` → **200 OK** or **404**

### Update user (name/email only)
`PATCH /users/{id}`
```json
{ "name": "Alice Doe", "email": "alice.doe@example.com" }
```
→ **200 OK**

### List users (optionally by group)
`GET /users?group=adult-1` → **200 OK**

## Grouping Rules
- 0–12 → `child`
- 13–17 → `teen`
- 18–64 → `adult`
- 65+ → `senior`
- Capacity per group: **3**. When full, the next numbered group is created (`adult-2`, `senior-3`, ...).

## Design Notes
- Group allocation happens **inside a DB transaction**. We select a row from `groups` where `member_count < capacity` using a row lock (`FOR UPDATE` via GORM), or create the next group with `index = MAX(index)+1`.
- The `group` field on `User` is **read-only** at the API level and always reflects computed assignment.

## Testing
```bash
go test ./...
```
(uses SQLite in-memory)

## Deployment
- Build image: `docker build -t abwaab/backend-task:latest .`
- Push to registry of your choice.
- AWS: run on ECS Fargate or EC2 + docker-compose. Provide env `DB_DSN` to point to your AWS RDS Postgres.

## Contact
Prepared for Abwaab backend task.
