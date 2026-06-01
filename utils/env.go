package utils

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"unicode"
)

func loadEnvFile() {
	if envFile, err := os.ReadFile(".env"); err != nil {
		slog.Debug(fmt.Sprintf("env file was not read %v", err))
	} else {
		for line := range bytes.Lines(envFile) {
			line = bytes.TrimRightFunc(line, unicode.IsControl)
			parts := bytes.SplitN(line, []byte("="), 2)
			if len(parts) != 2 {
				continue
			} else {
				name, value := string(parts[0]), string(parts[1])
				if _, exists := os.LookupEnv(name); !exists {
					os.Setenv(name, value)
				}
			}
		}
	}
}

func init() {
	loadEnvFile()
}

func Getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic(fmt.Sprintf("'%s' is not set or is unset", name))
	}

	return v
}

func GetenvAs[T any](name string, parse func(string) (T, error)) T {
	str := Getenv(name)
	if v, err := parse(str); err != nil {
		panic(fmt.Sprintf("failed to parse environment variable '%s' with value [%s]: %v", name, str, err))
	} else {
		return v
	}
}
