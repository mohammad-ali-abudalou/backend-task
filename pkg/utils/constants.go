package utils

// HTTP Status Codes :
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// Common API Error Messages :
const (
	ErrNoEnvFileFound                = "No .env File Found"
	ErrFailedConnectDatabase         = "Failed To Connect To Database"
	ErrDatabaseConnectedSuccessfully = "Database Connected Successfully"
	ErrUnsupportedDBDriver           = "Unsupported DB Driver"
	ErrMigrationFailed               = "Migration Failed"
	ErrInvalidRequestBody            = "Invalid Request Body"
	ErrInvalidId                     = "Invalid Id"
	ErrUserNotFound                  = "User Not Found"
	ErrNameIsRequired                = "Name Is Required"
	ErrInvalidEmailFormat            = "Invalid Email Format"
	ErrDateOfBirthFormat             = "date_of_birth Must Be YYYY-MM-DD"
	ErrEmailAlreadyExists            = "Email Already Exists"
	ErrNameCanNotEmpty               = "Name Cannot Be Empty"
	ErrRecordNotFound                = "Record Not Found"
	ErrInternalError                 = "Internal Server Error"
	ErrDateOfBirthCanNotInFuture     = "date_of_birth Cannot Be In The Future"
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
