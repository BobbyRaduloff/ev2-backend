package utils

import (
	"os"
	"strconv"
	"time"
)

func ParseEnvInt(envVar string, defaultValue int) int {
	valueStr := os.Getenv(envVar)
	if valueStr != "" {
		valueParsed, err := strconv.Atoi(valueStr)
		if err == nil {
			return valueParsed
		}
	}

	return defaultValue
}

func ParseEnvDuration(envVar string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(envVar)
	if valueStr != "" {
		valueParsed, err := strconv.Atoi(valueStr)
		if err == nil {
			return time.Duration(valueParsed) * time.Second
		}
	}

	return defaultValue
}

func ParseEnvString(envVar string, defaultValue string) string {
	valueStr := os.Getenv(envVar)
	if valueStr != "" {
		return valueStr
	}

	return defaultValue
}
