package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strconv"
)

type envelope map[string]any

func (app *application) writeJSON(w http.ResponseWriter, statusCode int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.WriteHeader(statusCode)

	w.Write(js)
	return nil
}

func getEnvInt(key string) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		slog.Info("error getting env variable", "key", key)
		panic(err)
	}

	return value
}
