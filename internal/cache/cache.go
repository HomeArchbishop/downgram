package cache

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

var cachePath string
var cacheData map[string]interface{}

func init() {
	executablePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executableDir := filepath.Dir(executablePath)

	cachePath = filepath.Join(executableDir, "CACHE")

	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		writeErr := os.WriteFile(cachePath, []byte("{}"), 0644)
		if writeErr != nil {
			panic(writeErr)
		}
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &cacheData)
	if err != nil {
		panic(err)
	}
}

func Get[T any](key string, defaultVal T) T {
	if val, ok := cacheData[key]; ok {
		exactVal, ok := val.(T)
		if ok {
			return exactVal
		} else {
			return defaultVal
		}
	}
	return defaultVal
}

func Update(key string, value any) {
	cacheData[key] = value
	data, err := json.Marshal(cacheData)
	if err != nil {
		panic(err)
	}

	writeErr := os.WriteFile(cachePath, data, 0644)
	if writeErr != nil {
		panic(writeErr)
	}
}
