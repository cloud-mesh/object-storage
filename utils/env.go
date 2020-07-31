package utils

import (
	"os"
	"strconv"
	"strings"
)

func Env(name string, defaultValue string) string {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		return defaultValue
	}

	return val
}

func EnvBool(name string, defaultValue bool) bool {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		return defaultValue
	}

	if b, err := strconv.ParseBool(val); err == nil {
		return b
	} else {
		return defaultValue
	}
}

func EnvInt(name string, defaultValue int) int {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		return defaultValue
	}

	if intVal, err := strconv.Atoi(val); err == nil {
		return intVal
	} else {
		return defaultValue
	}
}

func EnvArray(name string, sep string, defaultValue []string) []string {
	val, ok := os.LookupEnv(name)
	if !ok || val == "" {
		return defaultValue
	}

	arr := strings.Split(val, sep)
	resultArr := make([]string, 0, len(arr))
	for _, item := range arr {
		if strings.TrimSpace(item) != "" {
			resultArr = append(resultArr, item)
		}
	}
	return resultArr
}
