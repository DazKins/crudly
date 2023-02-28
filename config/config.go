package config

import (
	"crudly/util"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                  uint
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
		Port:                  getUint("PORT").UnwrapOrDefault(80),
		PostgresHost:          getEnv("POSTGRES_HOST").Unwrap(),
		PostgresPort:          getUint("POSTGRES_PORT").Unwrap(),
		PostgresUsername:      getEnv("POSTGRES_USERNAME").Unwrap(),
		PostgresPassword:      getEnv("POSTGRES_PASSWORD").Unwrap(),
		PostgresDatabase:      getEnv("POSTGRES_DATABASE").Unwrap(),
		ProjectCreationApiKey: getEnv("PROJECT_CREATION_API_KEY").Unwrap(),
	}
}

func getUint(env string) util.Result[uint] {
	numResult := getInt(env)

	if numResult.IsErr() {
		return util.ResultErr[uint](numResult.UnwrapErr())
	}

	num := numResult.Unwrap()

	if num < 0 {
		return util.ResultErr[uint](fmt.Errorf("env var: %d is less than zero", num))
	}

	return util.ResultOk(uint(num))
}

func getInt(env string) util.Result[int] {
	envResult := getEnv(env)

	if envResult.IsErr() {
		return util.ResultErr[int](envResult.UnwrapErr())
	}

	num, err := strconv.Atoi(envResult.Unwrap())

	if err != nil {
		return util.ResultErr[int](fmt.Errorf("env var: %s is not a valid number", envResult.Unwrap()))
	}

	return util.ResultOk(num)
}

func getEnv(env string) util.Result[string] {
	envVar, present := os.LookupEnv(env)

	if !present {
		return util.ResultErr[string](fmt.Errorf("env var: %s not present", env))
	}

	return util.ResultOk(envVar)
}
