package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	HttpsEnabled          bool
	SslCertificateFile    string
	SslPrivateKeyFile     string
	PostgresHost          string
	PostgresPort          uint
	PostgresUsername      string
	PostgresPassword      string
	PostgresDatabase      string
	ProjectCreationApiKey string
}

func InitialiseConfg() Config {
	godotenv.Load()

	return Config{
		HttpsEnabled:          getBoolOrPanic("HTTPS_ENABLED"),
		SslCertificateFile:    getEnvOrDefault("SSL_CERTIFICATE_FILE", ""),
		SslPrivateKeyFile:     getEnvOrDefault("SSL_PRIVATE_KEY_FILE", ""),
		PostgresHost:          getEnvOrPanic("POSTGRES_HOST"),
		PostgresPort:          getUintOrPanic("POSTGRES_PORT"),
		PostgresUsername:      getEnvOrPanic("POSTGRES_USERNAME"),
		PostgresPassword:      getEnvOrPanic("POSTGRES_PASSWORD"),
		PostgresDatabase:      getEnvOrPanic("POSTGRES_DATABASE"),
		ProjectCreationApiKey: getEnvOrPanic("PROJECT_CREATION_API_KEY"),
	}
}

func getBoolOrPanic(env string) bool {
	envVar := getEnvOrPanic(env)

	lower := strings.ToLower(envVar)

	if lower == "false" {
		return false
	}

	if lower == "true" {
		return true
	}

	panic(fmt.Sprintf("env var: %s is not a valid boolean", envVar))
}

func getUintOrPanic(env string) uint {
	num := getIntOrPanic(env)

	if num < 0 {
		panic(fmt.Sprintf("env var: %d is less than zero", num))
	}

	return uint(num)
}

func getIntOrPanic(env string) int {
	envVar := getEnvOrPanic(env)

	num, err := strconv.Atoi(envVar)

	if err != nil {
		panic(fmt.Sprintf("env var: %s is not a valid number", envVar))
	}

	return num
}

func getEnvOrPanic(env string) string {
	envVar, present := os.LookupEnv(env)

	if !present {
		panic(fmt.Sprintf("env var: %s not present", env))
	}

	return envVar
}

func getEnvOrDefault(env string, def string) string {
	envVar, present := os.LookupEnv(env)

	if !present {
		return def
	}

	return envVar
}
