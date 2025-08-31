package utils

//
const (
	AppName = "Backend TASK API"
	Version = "1.0.0"
)

// HTTP Status Codes :
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// Driver Type :
const (
	DriverPostgres = "postgres"
	DriverSqlite   = "sqlite"
)

// Age To Base Group ( Group Assignment ) :
const (
	BaseGroupChild  = "child"
	BaseGroupTeen   = "teen"
	BaseGroupAdult  = "adult"
	BaseGroupSenior = "senior"
	BaseGroupUnset  = "unset"
)
