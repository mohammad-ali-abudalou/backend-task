package config

import (
	"os"

	"backend-task/internal/utils"

	"github.com/joho/godotenv"
)

// LoadEnv Loads Environment Variables From A .env File ( If Available )
// And Falls Back To System Environment Variables When Running In Production.
func LoadEnv() {

	// Load .env Only If Present ( Don't Fail If Missing In Production ) :
	if err := godotenv.Load(); err != nil {
		utils.Info(utils.ErrNoEnvFileFound.Error())
	}
}

// GetEnv Retrieves An Environment Variable Or Returns A Fallback Value If Not Set.
func GetEnv(key, fallback string) string {

	if value, exists := os.LookupEnv(key); exists {

		return value
	}
	return fallback
}
