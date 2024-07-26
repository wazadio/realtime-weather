package pkg

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func LoadEnv(fileNames ...string) {
	fileName := ".env"
	if len(fileNames) > 0 {
		fileName = fileNames[0]
	}
	envFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("error loading env file : %s", err.Error())
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		if scanner.Text() == "" || string(scanner.Text()[0]) == "#" {
			continue
		}

		env := strings.SplitN(scanner.Text(), "=", 2)
		if len(env) <= 1 {
			continue
		}

		key := strings.TrimSpace(env[0])
		value := strings.TrimSpace(env[1])
		value = strings.ReplaceAll(value, "'", "")
		os.Setenv(key, value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
}
