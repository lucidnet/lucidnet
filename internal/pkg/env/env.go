package env

import "os"

func GetOrDefault(key string, def string) string {
	value := os.Getenv(key)

	if value == "" {
		return def
	}

	return value
}
