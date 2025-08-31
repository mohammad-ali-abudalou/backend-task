package config

import (
	"github.com/joho/godotenv"

	"backend-task/internal/utils"
)

func LoadEnv() {

	// Load Values From .env File.
	err := godotenv.Load()
	if err != nil {
		utils.Error(utils.ErrNoEnvFileFound.Error())
	}
}
