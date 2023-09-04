package env

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func Getenv(key string, defval string) string {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return defval
}
