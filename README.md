````markdown
# Backend Task (Go + Gin + GORM + Postgres)

Implements a user management REST API with automatic, capacity-aware age-group assignment.

## Features
- REST API (Gin)
- ORM (GORM) with Postgres (tests use SQLite)
- Concurrency-safe group allocation via dedicated `groups` table + row locks
- Validation (email format, uniqueness, DOB in the past)
- Swagger (OpenAPI 3) auto-generated docs at `/swagger/index.html`
- Unit tests
- Safe handling of concurrent user creation

---

## AWS Deployment

The application is deployed on an **AWS EC2 (Windows)** instance.  
The database uses **AWS RDS (PostgreSQL)** with public access enabled.  

- **Swagger API:** http://51.21.3.224:8080/swagger/index.html#/users/
- **RDS Connection Details:**  
  - Endpoint: `database-2.c1uuy8cm86e8.eu-north-1.rds.amazonaws.com`  
  - Port: `5432`  
  - Username: `postgres`  
  - Password: `Database`

You can connect and test the API directly using the above endpoint.

---

## Run Locally (No Docker)

1. Install and start **PostgreSQL** locally.
2. Set environment variables:

```bash
export DB_DRIVER=postgres
export DB_DSN="host=database-2.c1uuy8cm86e8.eu-north-1.rds.amazonaws.com user=postgres password=Database dbname=postgres port=5432 sslmode=require TimeZone=UTC"
````

3. Start the API:

```bash
go run ./cmd
```

> **Note for Windows (PowerShell):**
>
> * Use `go run ./cmd` (not `go run cmd/main.go`) to ensure module imports resolve correctly.
> * Swagger docs are auto-generated via:
>
> ```powershell
> swag init -g cmd/main.go
> ```

---

## Endpoints

### Create User

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

---

### Get User by ID

`GET /users/{id}` → **200 OK** or **404 Not Found**

---

### Update User (Name/Email Only)

`PATCH /users/{id}`

```json
{
  "name": "Alice Doe",
  "email": "alice.doe@example.com"
}
```

→ **200 OK**

---

### List Users (Optionally by Group)

`GET /users?group=adult-1` → **200 OK**
`GET /users` → List all users

---

## Grouping Rules

* 0 – 12 → `child`
* 13 – 17 → `teen`
* 18 – 64 → `adult`
* 65+ → `senior`
* Capacity per group: **3**. When full, the next numbered group is created (`adult-2`, `senior-3`, ... ).

---

## Design Notes

* Group allocation occurs **inside a DB transaction**.
* Rows from `groups` with `member_count < capacity` are selected using a **row lock** (`FOR UPDATE` via GORM).
  If full, a new group is created (`index = MAX(index)+1`).
* The `group` field on `User` is **read-only** at the API level.
* Swagger annotations (`@Summary`, `@Description`, `@Tags`, etc.) are included in handler functions.

---

## Testing

```bash
go test ./...
```

> Uses SQLite in-memory for testing.

---

## Deployment

### Prerequisites

* AWS Account with EC2 and RDS access
* Installed Go runtime on EC2 instance
* PostgreSQL RDS instance (with public access enabled)
* Security groups allowing inbound traffic on ports `8080` (API) and `5432` (Postgres)

### Steps

1. **Provision AWS RDS (PostgreSQL):**

   * Create an RDS PostgreSQL instance.
   * Enable public access.
   * Note down the `Endpoint`, `Port`, `Username`, and `Password`.

2. **Provision AWS EC2 Instance:**

   * Launch a Windows (or Linux) EC2 instance.
   * Open port `8080` in the Security Group.
   * SSH/RDP into the instance.

3. **Install Go and Git (on EC2):**

   ```bash
   choco install golang git -y   # Windows with Chocolatey
   ```

4. **Clone the Repository:**

   ```bash
   git clone https://github.com/mohammad-ali-abudalou/backend-task.git
   cd backend-task
   ```

5. **Configure Environment Variables (example for PowerShell):**

   ```powershell
   $env:DB_DRIVER="postgres"
   $env:DB_DSN="host=database-2.c1uuy8cm86e8.eu-north-1.rds.amazonaws.com user=postgres password=Database dbname=postgres port=5432 sslmode=disable TimeZone=UTC"
   ```

6. **Run the Application:**

   ```powershell
   go run ./cmd
   ```

7. **Verify Deployment:**

   * API available at: `http://<EC2-PUBLIC-IP>:8080` # http://51.21.3.224:8080/
   * Swagger available at: `http://<EC2-PUBLIC-IP>:8080/swagger/index.html` # http://51.21.3.224:8080/swagger/index.html#/users/

---

## Contact

For any questions or clarifications regarding this project:

**Mohammad Ali Abu-Dalou**

* Mobile: +962790132315
* Email: abudalou.mohammad@gmail.com
* LinkedIn: https://www.linkedin.com/in/mohammad-ali-abudalou/

```

---
