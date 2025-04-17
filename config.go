package main

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"sync"
)

type ConfigType struct {
	LogLevel slog.Level
	Debug    bool
	Secret   []byte
	Host     string
	Port     string
	DbUrl    string
}

const EnvPrefix string = "USER_AUTH"

var lock sync.Mutex
var config *ConfigType

func parseBool(val string) bool {
	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		panic(fmt.Sprintf("could not parse bool value: %s", val))
	}
	return boolVal
}

func parseLevel(val string) slog.Level {
	var level slog.Level
	err := level.UnmarshalText([]byte(val))
	if err != nil {
		panic(fmt.Sprintf("could not parse log level: %s", val))
	}
	return level
}

func parseSecret(val string) []byte {
	var secret []byte
	sz, err := base64.StdEncoding.Decode(secret, []byte(val))
	if err != nil || sz == 0 {
		panic("could not parse secret")
	}
	return secret
}

func getEnv(env string) string {
	var key = EnvPrefix + "_" + env
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("env var not defined: %s", key))
	}
	return val
}

// Get the singleton config
func NewConf() *ConfigType {
	if config != nil {
		return config
	}
	lock.Lock()
	defer lock.Unlock()
	config = &ConfigType{
		LogLevel: parseLevel(getEnv("LOG_LEVEL")),
		Debug:    parseBool(getEnv("DEBUG")),
		Secret:   parseSecret(getEnv("API_SECRET")),
		Host:     getEnv("HOST"),
		Port:     getEnv("PORT"),
		DbUrl:    getEnv("DB_URL"),
	}
	return config
}
