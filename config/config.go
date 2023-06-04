package config

import (
	"crudly/util/result"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             uint
	PostgresHost     string
	PostgresPort     uint
	PostgresUsername string
	PostgresPassword string
	PostgresDatabase string
	RedisHost        string
	AdminApiKey      string
}

func InitialiseConfg() Config {
	godotenv.Load()

	return Config{
		Port:             getUint("PORT").UnwrapOrDefault(80),
		PostgresHost:     getEnv("POSTGRES_HOST").Unwrap(),
		PostgresPort:     getUint("POSTGRES_PORT").Unwrap(),
		PostgresUsername: getEnv("POSTGRES_USERNAME").Unwrap(),
		PostgresPassword: getEnv("POSTGRES_PASSWORD").Unwrap(),
		PostgresDatabase: getEnv("POSTGRES_DATABASE").Unwrap(),
		RedisHost:        getEnv("REDIS_HOST").Unwrap(),
		AdminApiKey:      getEnv("ADMIN_API_KEY").Unwrap(),
	}
}

func getUint(env string) result.R[uint] {
	numResult := getInt(env)

	if numResult.IsErr() {
		return result.Err[uint](numResult.UnwrapErr())
	}

	num := numResult.Unwrap()

	if num < 0 {
		return result.Err[uint](fmt.Errorf("env var: %d is less than zero", num))
	}

	return result.Ok(uint(num))
}

func getInt(env string) result.R[int] {
	envResult := getEnv(env)

	if envResult.IsErr() {
		return result.Err[int](envResult.UnwrapErr())
	}

	num, err := strconv.Atoi(envResult.Unwrap())

	if err != nil {
		return result.Err[int](fmt.Errorf("env var: %s is not a valid number", envResult.Unwrap()))
	}

	return result.Ok(num)
}

func getEnv(env string) result.R[string] {
	envVar, present := os.LookupEnv(env)

	if !present {
		return result.Err[string](fmt.Errorf("env var: %s not present", env))
	}

	return result.Ok(envVar)
}
