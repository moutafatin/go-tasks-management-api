package main

import (
	"log/slog"
	"os"
	"strconv"
)

func getEnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		slog.Info("error getting env variable", "key", key)
		panic(err)
	}

	return value
}
