package util

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog/log"
)

// Read a file, exit if err
func ReadFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal().
			Str("path", path).
			Err(err).
			Msg("Failed to read file")
	}
	return b
}

// Write a file, exit if err
func WriteFile(path string, data []byte) {
	if err := ioutil.WriteFile(path, data, 0600); err != nil {
		log.Fatal().
			Str("path", path).
			Err(err).
			Msg("Failed to write file")
	}
}

// Check if a file exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
