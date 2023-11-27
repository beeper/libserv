package flagenv

import (
	"os"
	"strconv"
	"time"
)

func StringEnvWithDefault(key, defaultValue string) string {
	if envValue, ok := os.LookupEnv(key); !ok {
		return defaultValue
	} else {
		return envValue
	}
}

func BoolEnvWithDefault(key string, defaultValue bool) bool {
	if envValue, ok := os.LookupEnv(key); !ok {
		return defaultValue
	} else {
		return envValue == "true"
	}
}

func IntEnvWithDefault(key string, defaultValue int) int {
	if envValue, ok := os.LookupEnv(key); !ok {
		return defaultValue
	} else {
		v, err := strconv.Atoi(envValue)
		if err != nil {
			panic("invalid int env:" + key)
		}
		return v
	}
}

func DurationEnvWithDefault(key string, defaultValue time.Duration) time.Duration {
	if envValue, ok := os.LookupEnv(key); !ok {
		return defaultValue
	} else {
		v, err := time.ParseDuration(envValue)
		if err != nil {
			panic("invalid duation env:" + key)
		}
		return v
	}
}
