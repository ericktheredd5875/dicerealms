package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func init() {

	envFile := ""
	if IsTestEnv() {

		root, err := FindProjectRoot("")
		if err != nil {
			log.Fatalf("failed to find project root: %v", err)
		}

		envFile = filepath.Join(root, ".env.test")
		// log.Println("ℹ️ APP_ENV is not set, using .env.test")
	} else {
		envFile = ".env.dev"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Println("ℹ️ .env file not found")
	}

	// log.Println("🔑 Environment variables loaded from", appEnv)
	// appEnv := os.Getenv("APP_ENV")
	// log.Println("🔑 APP_ENV:", appEnv)
}

// ObtainEnv: Get an environment variable with a fallback
func ObtainEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}

// RequiredEnv: Check if required environment variables are set
func RequiredEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required environment variable: %s", key)
	}

	return val
}

func IsTestEnv() bool {
	is := false
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			is = true
			break
		}
	}

	return is
}
