package env

import (
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load(".env.local")
	godotenv.Load() // The Original .env
}