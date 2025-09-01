package utils

// ---------------- App Info ----------------
const (
	AppName = "Backend TASK API"
	Version = "1.0.0"
)

// ---------------- HTTP Status Codes ----------------
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// ---------------- Database Driver Types ----------------
const (
	DriverPostgres = "postgres"
	DriverSqlite   = "sqlite"
)

// ---------------- Age To Base Group ( For Group Assignment ) ----------------
const (
	BaseGroupChild  = "child"
	BaseGroupTeen   = "teen"
	BaseGroupAdult  = "adult"
	BaseGroupSenior = "senior"
	BaseGroupUnset  = "unset"
)

// ---------------- Data Source Name ( DSN ) Environment Keys ----------------
const (
	DSN_DRIVER_NAME = "DRIVER_NAME"
	DSN_DB_HOST     = "DB_HOST"
	DSN_DB_USER     = "DB_USER"
	DSN_DB_PASSWORD = "DB_PASSWORD"
	DSN_DB_NAME     = "DB_NAME"
	DSN_DB_PORT     = "DB_PORT"
	DSN_DB_SSLMODE  = "DB_SSLMODE"

	AUTO_MIGRATE = "AUTO_MIGRATE"
)

// ---------------- Group Settings ----------------
const (
	GroupCapacity = 3 // Maximum users per group
)
