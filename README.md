
# Backend Task (Go + Gin + GORM + Postgres)

Implements a **user management REST API** with automatic, capacity-aware **age-group assignment**.

---

## Features

- REST API with ( **Gin** ).
- Database layer using / ORM **GORM** with **Postgres** (tests run on mock).
- **Concurrency-safe** group allocation via a dedicated `groups` table + row-level locks.
- **Validation** :
  - Unique & valid email format.
  - DOB must be in the past.
- **Swagger (OpenAPI 3)** auto-generated docs at `/swagger/index.html`
- **Unit tests** with mocks.
- Safe handling of **concurrent user creation**.

---

## Live AWS Deployment

The application is currently deployed on :

- The application is deployed on an **AWS EC2 ( Windows )** instance.
- The database uses **AWS RDS (PostgreSQL)** with public access enabled.

**Live Swagger API Docs - Swagger API:**:  
http://51.21.3.224:8080/swagger/index.html#/

**Database Connection Details ( AWS RDS )**:  
- Endpoint: `database-2.c1uuy8cm86e8.eu-north-1.rds.amazonaws.com`  
- Port: `5432`  
- Username: `postgres`  
- Password: `Database`

You can connect and test the API directly using the above endpoint.

---


## Run Locally ( No AWS )

1. Install **PostgreSQL** locally and ensure it’s running.

2. Export environment variables:

```bash
export DB_DRIVER=postgres
export DB_DSN="host=database-2.c1uuy8cm86e8.eu-north-1.rds.amazonaws.com user=postgres password=Database dbname=postgres port=5432 sslmode=require TimeZone=UTC""
```

3. Start the API :

```bash
go run ./cmd/app
```

4. Swagger Docs :

```bash
swag init -g cmd/app/main.go
```
Then open: `http://localhost:8080/swagger/index.html`


> **Note for Windows (PowerShell):**
>
> * Use `go run ./cmd/app` to ensure module imports resolve correctly.
> * Swagger docs are auto-generated via:
>
> ```powershell
> swag init -g cmd/app/main.go
> ```

---


> **Windows ( PowerShell ) users :**  
> 
> ```powershell
> go run ./cmd/app
> swag init -g cmd/app/main.go   # Regenerate Swagger Docs.
> ```

---

## API Endpoints

### Create User

**`POST /users`**

**Request:**
```json
{
  "name": "Alice",
  "email": "alice@example.com",
  "date_of_birth": "1990-05-10"
}
```

**Response ( 201 Created ) :**

```json
{
  "id": "<uuid>",
  "name": "Alice",
  "email": "alice@example.com",
  "date_of_birth": "1990-05-10T00:00:00Z",
  "group": "adult-1",
  "created_at": "2025-09-01T10:05:00Z",
  "updated_at": "2025-09-01T10:05:00Z"
}
```

---

### Get User by ID
**GET /users/{id}**


- **200 OK** → User found.
- **404 Not Found** → Invalid ID.

---


### Update User ( Name / Email Only )
**PUT /users/{id}**

**Request:**
```json
{
  "name": "Alice Doe",
  "email": "alice.doe@example.com"
}
```

**Response ( 200 OK ) :**

```json
{
  "id": "<uuid>",
  "name": "Alice Doe",
  "email": "alice.doe@example.com",
  "date_of_birth": "1990-05-10T00:00:00Z",
  "group": "adult-1",
  "created_at": "2025-09-01T10:05:00Z",
  "updated_at": "2025-09-01T10:05:00Z"
}
```

---

### List Users ( Optionally By Group Filter )

**GET /users** → List all users  
**GET /users?group=adult-1** → List only users in `adult-1` 

**Response ( 200 OK ) :**
```json
[
  {
    "id": "8e0b2cfa-5df8-4c29-9b89-4e7a1fb3a411",
    "name": "Alice Doe",
    "email": "alice.doe@example.com",
    "date_of_birth": "1990-05-10T00:00:00Z",
    "group": "adult-1",
    "created_at": "2025-09-01T10:05:00Z",
    "updated_at": "2025-09-01T10:05:00Z"
  }
]
```

---

## Grouping Rules :

| Age Range | Group Name | Example |
|-----------|------------|---------|
| 0–12      | child      | `child-1`, `child-2` |
| 13–17     | teen       | `teen-1` |
| 18–64     | adult      | `adult-1`, `adult-2` |
| 65+       | senior     | `senior-1` |

- **Capacity per group:** 3  
- When full, the next numbered group is created (`adult-2`, `senior-3`, ...).

---

## Design Notes

* Group allocation occurs **inside a DB transaction**.
* Rows from `groups` with `member_count < capacity` are selected using a **row lock** (`FOR UPDATE` via GORM).
* If all groups are full, a new group is created automatically (`index = MAX(index)+1`).
* The `group` field on `User` is **read-only** at the API level.
* **Swagger annotations** (`@Summary`, `@Description`, `@Tags`, etc.) are included in all handler functions.

---

## Testing

Run unit tests with :

```bash
go test ./...
```

> Mocks are used for services and repositories..

---

## Deployment Guide ( AWS )

### Prerequisites

* AWS Account with EC2 and RDS access.
* Installed Go runtime on EC2 instance.
* PostgreSQL RDS instance ( with public access enabled ).
* Security groups allowing :
  - ports `8080` → API
  - `5432` → Postgres
- Installed Go & Git on EC2


### Steps

1. **Provision AWS RDS ( PostgreSQL ):**

   * Create an RDS PostgreSQL instance.
   * Enable public access.
   * Note down the `Endpoint`, `Port`, `Username`, and `Password`.

2. **Provision AWS EC2 Instance ( Windows / Linux ):**

   * Launch a Windows (or Linux) EC2 instance.
   * Open port `8080` in the Security Group.
   * SSH/RDP into the instance.

3. **Install Go and Git ( on EC2 ) :**

```bash
choco install golang git -y   # Windows With Chocolatey
```

4. **Clone the Repository :**

```bash
git clone https://github.com/mohammad-ali-abudalou/backend-task.git
cd backend-task
```

5. **Configure Environment Variables (example for PowerShell):**

```powershell
$env:DB_DRIVER="postgres"
$env:DB_DSN="host=database-2.c1uuy8cm86e8.eu-north-1.rds.amazonaws.com user=postgres password=Database dbname=postgres port=5432 sslmode=require TimeZone=UTC"
```

6. **Run the Application :**

```powershell
go run ./cmd/app
```

7. **Verify Deployment:**

   * API available at: `http://<EC2-PUBLIC-IP>:8080`
      ex. http://51.21.3.224:8080/
   * Swagger available at: `http://<EC2-PUBLIC-IP>:8080/swagger/index.html`
      ex. http://51.21.3.224:8080/swagger/index.html#/users/

---

## Contact

**Mohammad Ali Abu-Dalou**

* Mobile: +962790132315
* Email: abudalou.mohammad@gmail.com
* LinkedIn: https://www.linkedin.com/in/mohammad-ali-abudalou/

---