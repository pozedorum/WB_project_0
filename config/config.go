// Package credentials is responcible for checking accuracy of admin login and password
package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

var (
	path = "config/config.txt"

	errWrongParsingLine = errors.New("error with parsing credentials")
	errNotEnoughArgs    = errors.New("not enough credentials")
)

func GetDBConf() (result string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return "", errWrongParsingLine
	}
	in := bufio.NewScanner(file)
	in.Split(bufio.ScanLines)
	for in.Scan() {
		line := in.Text()
		credLine := strings.Split(line, "=")
		if credLine[0] == "DB_URL" {
			result = credLine[1]
		}
	}

	if result == "" {
		err = errNotEnoughArgs
	}
	return result, nil
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
