package yaml

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func IsValidEnvironmentName(
	name string,
) bool {
	if name == "" {
		return false
	}

	for _, character := range name {
		if unicode.IsLetter(character) ||
			unicode.IsDigit(character) ||
			character == '-' ||
			character == '_' {
			continue
		}

		return false
	}

	return true
}
func LoadEnvFile(
	path string,
	required bool,
) error {
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if required {
				return fmt.Errorf(
					"gorix config: selected environment file %q does not exist",
					path,
				)
			}

			// A default .env file is optional. Environment variables may have
			// been supplied by Docker, Kubernetes, CI, or the operating system.
			return nil
		}

		return fmt.Errorf(
			"gorix config: failed to open environment file %q: %w",
			path,
			err,
		)
	}
	defer file.Close()

	// Capture variables that existed before reading the dotenv file.
	// These values must never be overwritten.
	externalVariables := currentEnvironmentKeys()

	scanner := bufio.NewScanner(file)

	// Permit longer values such as certificates or large DSNs.
	scanner.Buffer(
		make([]byte, 1024),
		1024*1024,
	)

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())

		if line == "" ||
			strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(
				strings.TrimPrefix(line, "export "),
			)
		}

		key, rawValue, found := strings.Cut(
			line,
			"=",
		)
		if !found {
			return fmt.Errorf(
				"gorix config: invalid environment declaration in %s:%d",
				path,
				lineNumber,
			)
		}

		key = strings.TrimSpace(key)

		if !isValidEnvironmentKey(key) {
			return fmt.Errorf(
				"gorix config: invalid environment variable %q in %s:%d",
				key,
				path,
				lineNumber,
			)
		}

		value, err := parseEnvFileValue(rawValue)
		if err != nil {
			return fmt.Errorf(
				"gorix config: invalid value for %s in %s:%d: %w",
				key,
				path,
				lineNumber,
				err,
			)
		}

		// OS, Docker, Kubernetes, and CI values have higher priority.
		if _, protected := externalVariables[key]; protected {
			continue
		}

		// Allows one env value to refer to values declared earlier:
		//
		// DB_HOST=localhost
		// DB_DSN=postgres://user:pass@${DB_HOST}:5432/app
		value = os.Expand(
			value,
			func(name string) string {
				return os.Getenv(name)
			},
		)

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf(
				"gorix config: failed to set environment variable %q: %w",
				key,
				err,
			)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf(
			"gorix config: failed to read environment file %q: %w",
			path,
			err,
		)
	}

	return nil
}

func parseEnvFileValue(
	rawValue string,
) (string, error) {
	value := strings.TrimSpace(rawValue)

	if value == "" {
		return "", nil
	}

	if strings.HasPrefix(value, `"`) {
		if !strings.HasSuffix(value, `"`) {
			return "", fmt.Errorf(
				"unterminated double-quoted value",
			)
		}

		parsed, err := strconv.Unquote(value)
		if err != nil {
			return "", err
		}

		return parsed, nil
	}

	if strings.HasPrefix(value, "'") {
		if !strings.HasSuffix(value, "'") {
			return "", fmt.Errorf(
				"unterminated single-quoted value",
			)
		}

		return value[1 : len(value)-1], nil
	}

	// Supports:
	//
	// DB_PORT=5432 # PostgreSQL port
	//
	// A # without preceding whitespace remains part of the value.
	if index := strings.Index(value, " #"); index >= 0 {
		value = strings.TrimSpace(
			value[:index],
		)
	}

	return value, nil
}

func currentEnvironmentKeys() map[string]struct{} {
	result := make(
		map[string]struct{},
	)

	for _, entry := range os.Environ() {
		key, _, found := strings.Cut(
			entry,
			"=",
		)
		if found {
			result[key] = struct{}{}
		}
	}

	return result
}

func isValidEnvironmentKey(
	key string,
) bool {
	if key == "" {
		return false
	}

	for index, character := range key {
		if index == 0 {
			if character == '_' ||
				unicode.IsLetter(character) {
				continue
			}

			return false
		}

		if character == '_' ||
			unicode.IsLetter(character) ||
			unicode.IsDigit(character) {
			continue
		}

		return false
	}

	return true
}
