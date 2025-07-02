// Package credentials is responcible for checking accuracy of admin login and password
package config

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

var (
	path = "/app/config/config.txt"

	errWrongParsingLine = errors.New("error with parsing credentials")
	errNotEnoughArgs    = errors.New("not enough credentials")
)

func GetDBConf() (string, error) {

	log.Printf("Trying to read config from: %s", path)

	if url := os.Getenv("DB_URL"); url != "" {
		return url, nil
	}

	content, err := os.ReadFile(path)
	if err == nil {
		return parseConfig(string(content))
	}

	log.Printf("Config file read error")

	return "", errors.New("DB config not found")
}

func GetKafkaConf() (result string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errWrongParsingLine
	}
	in := bufio.NewScanner(file)
	in.Split(bufio.ScanLines)
	for in.Scan() {
		line := in.Text()
		credLine := strings.Split(line, "=")
		if credLine[0] == "KAFKA_BROKERS" {
			result = credLine[1]
		}
	}

	if result == "" {
		err = errNotEnoughArgs
	}
	return result, nil
}

func parseConfig(content string) (string, error) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "DB_URL=") {
			dbURL := strings.TrimPrefix(line, "DB_URL=")
			log.Printf("Found DB_URL: %s", dbURL)
			return strings.TrimSpace(dbURL), nil
		}
	}
	return "", errors.New("database URL not found")
}
