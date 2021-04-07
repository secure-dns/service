package env

import "github.com/joho/godotenv"

//Load environment variables from .env files.
func Load() {
	godotenv.Load(".env.local", ".env")
}
